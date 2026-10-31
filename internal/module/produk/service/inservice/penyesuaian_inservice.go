package inservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"

	"pos-jamu-backend/internal/database"
	jurnaldomain "pos-jamu-backend/internal/module/jurnal/tck/domain"
	jurnaltck "pos-jamu-backend/internal/module/jurnal/tck/inservice"
	"pos-jamu-backend/internal/module/produk/service/repository"
	"pos-jamu-backend/internal/module/produk/tck/domain"
	tck "pos-jamu-backend/internal/module/produk/tck/inservice"
	"pos-jamu-backend/internal/util"
	"pos-jamu-backend/internal/util/apperror"
)

type PenyesuaianStokInService struct {
	db     *sql.DB
	jurnal jurnaltck.IJurnalInService
}

var _ tck.IPenyesuaianStokInService = (*PenyesuaianStokInService)(nil)

func NewPenyesuaianStokInService(db *sql.DB, jurnal jurnaltck.IJurnalInService) *PenyesuaianStokInService {
	return &PenyesuaianStokInService{db: db, jurnal: jurnal}
}

// rencana = hasil hitung penyesuaian sebelum ditulis ke database.
type rencana struct {
	qty      int64
	hpp      decimal.Decimal // nilai per unit yang dicatat
	hppBaru  decimal.Decimal // produk.hpp_rata_rata setelah penyesuaian
	stokBaru int64
	baris    []jurnaldomain.Baris
}

// HitungPenyesuaian berisi aturan bisnis penyesuaian stok (tanpa database).
//   - stok_awal: stok += qty, HPP ikut rumus moving average; jurnal Persediaan / Modal Pemilik
//   - opname:    stok = stok fisik, HPP tetap; selisih dinilai dengan HPP saat ini
//     (berkurang: HPP / Persediaan, bertambah: Persediaan / HPP)
//   - rusak:     stok -= qty, HPP tetap; jurnal HPP / Persediaan
func HitungPenyesuaian(in domain.PenyesuaianInput, stok int64, hpp decimal.Decimal) (rencana, error) {
	switch in.Jenis {
	case domain.JenisStokAwal:
		if in.Qty <= 0 {
			return rencana{}, apperror.BadRequest("Jumlah stok awal harus lebih dari 0")
		}
		if in.HargaModal.IsNegative() {
			return rencana{}, apperror.BadRequest("Harga modal tidak boleh minus")
		}
		nilai := util.Round2(in.HargaModal.Mul(decimal.NewFromInt(in.Qty)))
		return rencana{
			qty:      in.Qty,
			hpp:      util.Round2(in.HargaModal),
			hppBaru:  util.HPPRataRataBaru(stok, hpp, in.Qty, in.HargaModal),
			stokBaru: stok + in.Qty,
			baris: []jurnaldomain.Baris{
				jurnaldomain.Debit(jurnaldomain.AkunPersediaan, nilai),
				jurnaldomain.Kredit(jurnaldomain.AkunModalPemilik, nilai),
			},
		}, nil

	case domain.JenisOpname:
		if in.StokFisik == nil || *in.StokFisik < 0 {
			return rencana{}, apperror.BadRequest("Stok fisik wajib diisi dan tidak boleh minus")
		}
		qty := *in.StokFisik - stok
		if qty == 0 {
			return rencana{}, apperror.BadRequest("Stok fisik sama dengan stok di sistem, tidak ada yang perlu disesuaikan")
		}
		nilai := util.Round2(hpp.Mul(decimal.NewFromInt(abs(qty))))
		baris := []jurnaldomain.Baris{
			jurnaldomain.Debit(jurnaldomain.AkunHPP, nilai),
			jurnaldomain.Kredit(jurnaldomain.AkunPersediaan, nilai),
		}
		if qty > 0 {
			baris = []jurnaldomain.Baris{
				jurnaldomain.Debit(jurnaldomain.AkunPersediaan, nilai),
				jurnaldomain.Kredit(jurnaldomain.AkunHPP, nilai),
			}
		}
		return rencana{qty: qty, hpp: hpp, hppBaru: hpp, stokBaru: *in.StokFisik, baris: baris}, nil

	case domain.JenisRusak:
		if in.Qty <= 0 {
			return rencana{}, apperror.BadRequest("Jumlah barang rusak harus lebih dari 0")
		}
		nilai := util.Round2(hpp.Mul(decimal.NewFromInt(in.Qty)))
		return rencana{
			qty:      -in.Qty,
			hpp:      hpp,
			hppBaru:  hpp,
			stokBaru: stok - in.Qty,
			baris: []jurnaldomain.Baris{
				jurnaldomain.Debit(jurnaldomain.AkunHPP, nilai),
				jurnaldomain.Kredit(jurnaldomain.AkunPersediaan, nilai),
			},
		}, nil
	}
	return rencana{}, apperror.BadRequest("Jenis penyesuaian harus stok_awal, opname, atau rusak")
}

func (s *PenyesuaianStokInService) Create(ctx context.Context, in domain.PenyesuaianInput) (domain.PenyesuaianStok, error) {
	var id uint64
	err := database.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		produkRepo := repository.NewProdukRepository(tx)

		// Kunci baris produk supaya tidak balapan dengan transaksi lain.
		cur, err := produkRepo.GetStokForUpdate(ctx, in.ProdukID)
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound("Produk tidak ditemukan")
		}
		if err != nil {
			return err
		}

		r, err := HitungPenyesuaian(in, cur.Stok, cur.HPPRataRata)
		if err != nil {
			return err
		}

		if err := produkRepo.UpdateStokHPP(ctx, in.ProdukID, r.stokBaru, r.hppBaru); err != nil {
			return err
		}

		var ket *string
		if in.Keterangan != "" {
			ket = &in.Keterangan
		}
		id, err = repository.NewPenyesuaianRepository(tx).Create(ctx, domain.PenyesuaianStok{
			UserID:      in.UserID,
			ProdukID:    in.ProdukID,
			Jenis:       in.Jenis,
			Qty:         r.qty,
			StokSebelum: cur.Stok,
			HPP:         r.hpp,
			Keterangan:  ket,
		})
		if err != nil {
			return err
		}

		_, err = s.jurnal.Post(ctx, tx, jurnaldomain.Entri{
			UserID:     in.UserID,
			Keterangan: fmt.Sprintf("Penyesuaian stok (%s) produk #%d", in.Jenis, in.ProdukID),
			SumberTipe: jurnaldomain.SumberPenyesuaianStok,
			SumberID:   id,
			Baris:      r.baris,
		})
		return err
	})
	if err != nil {
		return domain.PenyesuaianStok{}, asAppError(err)
	}

	p, err := repository.NewPenyesuaianRepository(s.db).GetByID(ctx, id)
	if err != nil {
		return domain.PenyesuaianStok{}, apperror.Internal(err)
	}
	return p, nil
}

func (s *PenyesuaianStokInService) ListByProduk(ctx context.Context, produkID uint64, limit int32) ([]domain.PenyesuaianStok, error) {
	list, err := repository.NewPenyesuaianRepository(s.db).ListByProduk(ctx, produkID, limit)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return list, nil
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
