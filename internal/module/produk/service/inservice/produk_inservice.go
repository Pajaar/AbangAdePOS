package inservice

import (
	"context"
	"database/sql"
	"errors"

	"pos-jamu-backend/internal/database"
	"pos-jamu-backend/internal/module/produk/service/repository"
	"pos-jamu-backend/internal/module/produk/tck/domain"
	tck "pos-jamu-backend/internal/module/produk/tck/inservice"
	"pos-jamu-backend/internal/util/apperror"
)

type ProdukInService struct {
	db *sql.DB
}

var _ tck.IProdukInService = (*ProdukInService)(nil)

func NewProdukInService(db *sql.DB) *ProdukInService {
	return &ProdukInService{db: db}
}

func (s *ProdukInService) List(ctx context.Context, f domain.ProdukFilter) ([]domain.Produk, int64, error) {
	repo := repository.NewProdukRepository(s.db)
	items, err := repo.List(ctx, f)
	if err != nil {
		return nil, 0, apperror.Internal(err)
	}
	total, err := repo.Count(ctx, f)
	if err != nil {
		return nil, 0, apperror.Internal(err)
	}
	return items, total, nil
}

func (s *ProdukInService) GetByID(ctx context.Context, id uint64) (domain.Produk, error) {
	p, err := repository.NewProdukRepository(s.db).GetByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Produk{}, apperror.NotFound("Produk tidak ditemukan")
	}
	if err != nil {
		return domain.Produk{}, apperror.Internal(err)
	}
	p.Barcodes, err = repository.NewBarcodeRepository(s.db).ListByProduk(ctx, id)
	if err != nil {
		return domain.Produk{}, apperror.Internal(err)
	}
	return p, nil
}

func (s *ProdukInService) Create(ctx context.Context, in domain.ProdukInput, barcodes []domain.BarcodeInput) (uint64, error) {
	// Normalisasi semua barcode dulu supaya input salah ditolak sebelum
	// transaksi dibuka.
	prepared := make([]domain.Barcode, 0, len(barcodes))
	seen := map[string]bool{}
	for _, b := range barcodes {
		nb, err := prepareBarcode(b)
		if err != nil {
			return 0, err
		}
		if seen[nb.Kode] {
			return 0, apperror.BadRequest("Barcode " + nb.Kode + " diisi lebih dari sekali")
		}
		seen[nb.Kode] = true
		prepared = append(prepared, nb)
	}

	var id uint64
	err := database.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		var err error
		id, err = repository.NewProdukRepository(tx).Create(ctx, in, domain.StatusDariHarga(in.HargaNormal))
		if err != nil {
			return mapProdukWriteErr(err)
		}
		barcodeRepo := repository.NewBarcodeRepository(tx)
		for _, b := range prepared {
			b.ProdukID = id
			if _, err := barcodeRepo.Create(ctx, b); err != nil {
				return mapBarcodeWriteErr(err, b.Kode)
			}
		}
		return nil
	})
	if err != nil {
		return 0, asAppError(err)
	}
	return id, nil
}

func (s *ProdukInService) Update(ctx context.Context, id uint64, in domain.ProdukInput) error {
	repo := repository.NewProdukRepository(s.db)
	if _, err := repo.GetByID(ctx, id); errors.Is(err, sql.ErrNoRows) {
		return apperror.NotFound("Produk tidak ditemukan")
	} else if err != nil {
		return apperror.Internal(err)
	}
	if err := repo.Update(ctx, id, in, domain.StatusDariHarga(in.HargaNormal)); err != nil {
		return asAppError(mapProdukWriteErr(err))
	}
	return nil
}

func mapProdukWriteErr(err error) error {
	if mysqlErrNo(err) == mysqlErrForeignKey {
		return apperror.BadRequest("Kategori tidak ditemukan")
	}
	return err
}
