// Package validator berisi aturan bisnis menu jamu yang dicek sebelum simpan.
package validator

import (
	"context"
	"errors"
	"net/http"

	"github.com/shopspring/decimal"

	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	produktck "pos-jamu-backend/internal/module/produk/tck/inservice"
	"pos-jamu-backend/internal/util/apperror"
)

// cekProduk: produk ada dan masih aktif. Dipakai jamu (sachet wajib) dan
// item tambahan (opsional).
func cekProduk(ctx context.Context, svc produktck.IProdukInService, id uint64) ([]processor.ValidationError, error) {
	p, err := svc.GetByID(ctx, id)
	var appErr *apperror.AppError
	if errors.As(err, &appErr) && appErr.Status == http.StatusNotFound {
		return processor.OneError("produk_id", "not_found", "Produk tidak ditemukan"), nil
	}
	if err != nil {
		return nil, err
	}
	if !p.Aktif {
		return processor.OneError("produk_id", "inactive", "Produk sudah tidak aktif"), nil
	}
	return processor.NoErrors(), nil
}

func cekHarga(harga decimal.Decimal) []processor.ValidationError {
	if !harga.IsPositive() {
		return processor.OneError("harga_jual", "min", "Harga jual harus lebih dari 0")
	}
	return nil
}

// JamuValidator: harga seduh > 0 dan produk sachet ada serta aktif.
type JamuValidator struct {
	produk produktck.IProdukInService
}

var _ processor.IValidatorProcessor[dto.JamuRequest] = (*JamuValidator)(nil)

func NewJamuValidator(p produktck.IProdukInService) *JamuValidator {
	return &JamuValidator{produk: p}
}

func (v *JamuValidator) Process(ctx context.Context, req dto.JamuRequest) ([]processor.ValidationError, error) {
	errs := cekHarga(req.HargaJual)
	produkErrs, err := cekProduk(ctx, v.produk, req.ProdukID)
	if err != nil {
		return nil, err
	}
	return append(errs, produkErrs...), nil
}

// AdditionalValidator: harga > 0; bila produk_id diisi, produk ada dan aktif.
type AdditionalValidator struct {
	produk produktck.IProdukInService
}

var _ processor.IValidatorProcessor[dto.AdditionalRequest] = (*AdditionalValidator)(nil)

func NewAdditionalValidator(p produktck.IProdukInService) *AdditionalValidator {
	return &AdditionalValidator{produk: p}
}

func (v *AdditionalValidator) Process(ctx context.Context, req dto.AdditionalRequest) ([]processor.ValidationError, error) {
	errs := cekHarga(req.HargaJual)
	if req.ProdukID == nil {
		return errs, nil
	}
	produkErrs, err := cekProduk(ctx, v.produk, *req.ProdukID)
	if err != nil {
		return nil, err
	}
	return append(errs, produkErrs...), nil
}
