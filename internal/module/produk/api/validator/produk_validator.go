package validator

import (
	"context"
	"errors"
	"net/http"

	"github.com/shopspring/decimal"

	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/produk/tck/inservice"
	"pos-jamu-backend/internal/util/apperror"
)

// ProdukKategoriValidator: kategori ada dan masih aktif.
type ProdukKategoriValidator struct {
	kategori inservice.IKategoriInService
}

var _ processor.IValidatorProcessor[dto.ProdukRequest] = (*ProdukKategoriValidator)(nil)

func NewProdukKategoriValidator(k inservice.IKategoriInService) *ProdukKategoriValidator {
	return &ProdukKategoriValidator{kategori: k}
}

func (v *ProdukKategoriValidator) Process(ctx context.Context, req dto.ProdukRequest) ([]processor.ValidationError, error) {
	k, err := v.kategori.GetByID(ctx, req.KategoriID)
	var appErr *apperror.AppError
	if errors.As(err, &appErr) && appErr.Status == http.StatusNotFound {
		return processor.OneError("kategori_id", "not_found", "Kategori tidak ditemukan"), nil
	}
	if err != nil {
		return nil, err
	}
	if !k.Aktif {
		return processor.OneError("kategori_id", "inactive", "Kategori sudah tidak aktif"), nil
	}
	return processor.NoErrors(), nil
}

var maxMargin = decimal.NewFromInt(1000)

// ProdukHargaValidator: harga tidak minus; produk yang bisa ditawar wajib
// punya margin minimum (dasar harga nett).
type ProdukHargaValidator struct{}

var _ processor.IValidatorProcessor[dto.ProdukRequest] = ProdukHargaValidator{}

func (ProdukHargaValidator) Process(_ context.Context, req dto.ProdukRequest) ([]processor.ValidationError, error) {
	var errs []processor.ValidationError
	if req.HargaNormal.IsNegative() {
		errs = append(errs, processor.ValidationError{Field: "harga_normal", Code: "min", Message: "Harga normal tidak boleh minus"})
	}
	if req.MarginMinPersen != nil && (req.MarginMinPersen.IsNegative() || req.MarginMinPersen.GreaterThan(maxMargin)) {
		errs = append(errs, processor.ValidationError{Field: "margin_min_persen", Code: "range", Message: "Margin minimum harus 0 sampai 1000 persen"})
	}
	if req.BisaDitawar && req.MarginMinPersen == nil {
		errs = append(errs, processor.ValidationError{Field: "margin_min_persen", Code: "required", Message: "Margin minimum wajib diisi untuk produk yang bisa ditawar"})
	}
	return errs, nil
}
