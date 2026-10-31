// Package inservice berisi implementasi posting jurnal otomatis.
package inservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"

	"pos-jamu-backend/internal/module/jurnal/service/repository"
	"pos-jamu-backend/internal/module/jurnal/tck/domain"
	tck "pos-jamu-backend/internal/module/jurnal/tck/inservice"
)

type JurnalInService struct{}

var _ tck.IJurnalInService = (*JurnalInService)(nil)

func NewJurnalInService() *JurnalInService {
	return &JurnalInService{}
}

func (s *JurnalInService) Post(ctx context.Context, tx *sql.Tx, e domain.Entri) (uint64, error) {
	baris, err := ValidasiBaris(e.Baris)
	if err != nil {
		return 0, err
	}
	if len(baris) == 0 {
		return 0, nil
	}

	repo := repository.NewJurnalRepository(tx)
	jurnalID, err := repo.CreateJurnal(ctx, e.UserID, e.Keterangan, e.SumberTipe, e.SumberID)
	if err != nil {
		return 0, err
	}
	for _, b := range baris {
		akunID, err := repo.GetAkunIDByKode(ctx, b.AkunKode)
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("akun %s tidak ada di tabel akun", b.AkunKode)
		}
		if err != nil {
			return 0, err
		}
		if err := repo.CreateDetail(ctx, jurnalID, akunID, b.Debit, b.Kredit); err != nil {
			return 0, err
		}
	}
	return jurnalID, nil
}

// ValidasiBaris membuang baris bernilai nol, lalu memastikan setiap baris
// hanya berisi satu sisi bernilai positif dan total debit = total kredit.
func ValidasiBaris(in []domain.Baris) ([]domain.Baris, error) {
	out := make([]domain.Baris, 0, len(in))
	totalDebit, totalKredit := decimal.Zero, decimal.Zero
	for _, b := range in {
		if b.Debit.IsZero() && b.Kredit.IsZero() {
			continue
		}
		if b.Debit.IsNegative() || b.Kredit.IsNegative() || (b.Debit.IsPositive() && b.Kredit.IsPositive()) {
			return nil, fmt.Errorf("baris jurnal akun %s tidak valid (debit %s, kredit %s)", b.AkunKode, b.Debit, b.Kredit)
		}
		totalDebit = totalDebit.Add(b.Debit)
		totalKredit = totalKredit.Add(b.Kredit)
		out = append(out, b)
	}
	if !totalDebit.Equal(totalKredit) {
		return nil, fmt.Errorf("jurnal tidak seimbang: debit %s, kredit %s", totalDebit, totalKredit)
	}
	return out, nil
}
