// Package validator berisi satu validator per aturan bisnis modul produk
// (setara api/processor/validator di primecash). Validator mengembalikan
// daftar ValidationError; error kedua hanya untuk kegagalan teknis.
package validator

import (
	"context"
	"strings"

	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/produk/tck/inservice"
)

// KategoriNamaUnikValidator: nama kategori belum dipakai kategori lain.
type KategoriNamaUnikValidator struct {
	kategori inservice.IKategoriInService
}

var _ processor.IValidatorProcessor[dto.KategoriRequest] = (*KategoriNamaUnikValidator)(nil)

func NewKategoriNamaUnikValidator(k inservice.IKategoriInService) *KategoriNamaUnikValidator {
	return &KategoriNamaUnikValidator{kategori: k}
}

func (v *KategoriNamaUnikValidator) Process(ctx context.Context, req dto.KategoriRequest) ([]processor.ValidationError, error) {
	k, found, err := v.kategori.FindByNama(ctx, strings.TrimSpace(req.Nama))
	if err != nil {
		return nil, err
	}
	if found && k.ID != req.ID {
		return processor.OneError("nama", "already_exists", "Nama kategori sudah dipakai"), nil
	}
	return processor.NoErrors(), nil
}
