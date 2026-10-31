package inservice

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"pos-jamu-backend/internal/module/produk/service/repository"
	"pos-jamu-backend/internal/module/produk/tck/domain"
	tck "pos-jamu-backend/internal/module/produk/tck/inservice"
	"pos-jamu-backend/internal/util"
	"pos-jamu-backend/internal/util/apperror"
)

const (
	sumberManual = "manual"
)

type BarcodeInService struct {
	db      *sql.DB
	produk  tck.IProdukInService
	sources []tck.IBarcodeSource // dicoba berurutan
}

var _ tck.IBarcodeInService = (*BarcodeInService)(nil)

func NewBarcodeInService(db *sql.DB, produk tck.IProdukInService, sources ...tck.IBarcodeSource) *BarcodeInService {
	return &BarcodeInService{db: db, produk: produk, sources: sources}
}

func (s *BarcodeInService) Add(ctx context.Context, produkID uint64, in domain.BarcodeInput) (domain.Barcode, error) {
	b, err := prepareBarcode(in)
	if err != nil {
		return domain.Barcode{}, err
	}
	if _, err := s.produk.GetByID(ctx, produkID); err != nil {
		return domain.Barcode{}, err
	}
	b.ProdukID = produkID

	repo := repository.NewBarcodeRepository(s.db)
	id, err := repo.Create(ctx, b)
	if err != nil {
		return domain.Barcode{}, asAppError(mapBarcodeWriteErr(err, b.Kode))
	}
	created, err := repo.GetByID(ctx, id)
	if err != nil {
		return domain.Barcode{}, apperror.Internal(err)
	}
	return created, nil
}

func (s *BarcodeInService) Delete(ctx context.Context, id uint64) error {
	repo := repository.NewBarcodeRepository(s.db)
	if _, err := repo.GetByID(ctx, id); errors.Is(err, sql.ErrNoRows) {
		return apperror.NotFound("Barcode tidak ditemukan")
	} else if err != nil {
		return apperror.Internal(err)
	}
	if err := repo.Delete(ctx, id); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (s *BarcodeInService) FindByKode(ctx context.Context, kode string) (domain.Barcode, bool, error) {
	b, err := repository.NewBarcodeRepository(s.db).GetByKode(ctx, kode)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Barcode{}, false, nil
	}
	if err != nil {
		return domain.Barcode{}, false, apperror.Internal(err)
	}
	return b, true, nil
}

func (s *BarcodeInService) FindForKasir(ctx context.Context, raw string) (domain.Produk, domain.Barcode, error) {
	kode, _, err := util.NormalizeBarcode(raw)
	if err != nil {
		return domain.Produk{}, domain.Barcode{}, apperror.BadRequest(err.Error())
	}
	b, found, err := s.FindByKode(ctx, kode)
	if err != nil {
		return domain.Produk{}, domain.Barcode{}, err
	}
	if !found {
		return domain.Produk{}, domain.Barcode{}, apperror.NotFound("Barcode belum terdaftar")
	}
	p, err := s.produk.GetByID(ctx, b.ProdukID)
	if err != nil {
		return domain.Produk{}, domain.Barcode{}, err
	}
	if !p.Aktif {
		return domain.Produk{}, domain.Barcode{}, apperror.NotFound("Produk sudah tidak aktif")
	}
	if p.Status == domain.StatusDraft {
		return domain.Produk{}, domain.Barcode{}, apperror.Conflict("Produk masih draft, lengkapi harga dulu di menu Produk")
	}
	return p, b, nil
}

// Lookup: lokal -> cache -> API publik berurutan. Hasil API (termasuk
// "tidak ditemukan" yang pasti) disimpan ke cache; error/timeout tidak.
func (s *BarcodeInService) Lookup(ctx context.Context, raw string) (domain.LookupResult, error) {
	kode, jenis, err := util.NormalizeBarcode(raw)
	if err != nil {
		return domain.LookupResult{}, apperror.BadRequest(err.Error())
	}
	res := domain.LookupResult{Kode: kode, Jenis: jenis}

	// 1. Sudah terdaftar di toko.
	b, found, err := s.FindByKode(ctx, kode)
	if err != nil {
		return res, err
	}
	if found {
		p, err := s.produk.GetByID(ctx, b.ProdukID)
		if err != nil {
			return res, err
		}
		res.Ditemukan, res.Sumber, res.Produk = true, domain.SumberLokal, &p
		res.Nama = &p.Nama
		return res, nil
	}

	// Kode internal toko tidak mungkin ada di database publik.
	if jenis == util.BarcodeInternal {
		return res, nil
	}

	repo := repository.NewBarcodeRepository(s.db)

	// 2. Cache.
	cached, err := repo.GetCache(ctx, kode)
	if err == nil {
		res.Ditemukan, res.Sumber, res.DariCache = cached.Ditemukan, cached.Sumber, true
		res.Nama, res.Merek = cached.Nama, cached.Merek
		return res, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return res, apperror.Internal(err)
	}

	// 3. API publik.
	gagal := false
	for _, src := range s.sources {
		ext, err := src.Lookup(ctx, kode)
		if err != nil {
			gagal = true
			slog.WarnContext(ctx, "barcode_lookup_gagal", "sumber", src.Name(), "barcode", kode, "error", err.Error())
			continue
		}
		if !ext.Found {
			continue
		}
		res.Ditemukan, res.Sumber = true, src.Name()
		res.Nama, res.Merek = nonEmpty(ext.Nama), nonEmpty(ext.Merek)
		s.saveCache(ctx, repo, kode, repository.CachedLookup{
			Sumber: src.Name(), Ditemukan: true, Nama: res.Nama, Merek: res.Merek,
		}, ext.Raw)
		return res, nil
	}

	res.SebagianGagal = gagal
	if !gagal && len(s.sources) > 0 {
		res.Sumber = domain.SumberSemua
		s.saveCache(ctx, repo, kode, repository.CachedLookup{Sumber: domain.SumberSemua, Ditemukan: false}, nil)
	}
	return res, nil
}

// saveCache tidak menggagalkan lookup bila penyimpanan cache error.
func (s *BarcodeInService) saveCache(ctx context.Context, repo repository.IBarcodeRepository, kode string, c repository.CachedLookup, raw []byte) {
	if err := repo.UpsertCache(ctx, kode, c, raw); err != nil {
		slog.ErrorContext(ctx, "barcode_cache_gagal", "barcode", kode, "error", err.Error())
	}
}

// prepareBarcode menormalisasi input barcode dan mengisi nilai default.
func prepareBarcode(in domain.BarcodeInput) (domain.Barcode, error) {
	kode, jenis, err := util.NormalizeBarcode(in.Kode)
	if err != nil {
		return domain.Barcode{}, apperror.BadRequest(err.Error())
	}
	isi := in.IsiPerKemasan
	if isi == 0 {
		isi = 1
	}
	if isi < 0 {
		return domain.Barcode{}, apperror.BadRequest("Isi per kemasan harus lebih dari 0")
	}
	sumber := in.Sumber
	if sumber == "" {
		sumber = sumberManual
	}
	return domain.Barcode{Kode: kode, Jenis: jenis, IsiPerKemasan: isi, Sumber: sumber}, nil
}

func mapBarcodeWriteErr(err error, kode string) error {
	if mysqlErrNo(err) == mysqlErrDuplicate {
		return apperror.Conflict("Barcode " + kode + " sudah dipakai produk lain")
	}
	return err
}

func nonEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
