package validator

import (
	"context"
	"fmt"

	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/produk/tck/inservice"
	"pos-jamu-backend/internal/util"
)

// cekBarcode: format barcode valid, tidak dobel di input, dan belum dipakai produk lain.
func cekBarcode(ctx context.Context, svc inservice.IBarcodeInService, field string, kodes []string) ([]processor.ValidationError, error) {
	var errs []processor.ValidationError
	seen := map[string]bool{}
	for i, raw := range kodes {
		f := field
		if len(kodes) > 1 || field != "kode" {
			f = fmt.Sprintf("%s[%d].kode", field, i)
		}
		kode, _, err := util.NormalizeBarcode(raw)
		if err != nil {
			errs = append(errs, processor.ValidationError{Field: f, Code: "invalid", Message: err.Error()})
			continue
		}
		if seen[kode] {
			errs = append(errs, processor.ValidationError{Field: f, Code: "duplicate", Message: "Barcode " + kode + " diisi lebih dari sekali"})
			continue
		}
		seen[kode] = true
		if _, found, err := svc.FindByKode(ctx, kode); err != nil {
			return nil, err
		} else if found {
			errs = append(errs, processor.ValidationError{Field: f, Code: "already_exists", Message: "Barcode " + kode + " sudah dipakai produk lain"})
		}
	}
	return errs, nil
}

// ProdukBarcodeBelumDipakaiValidator: barcode di form tambah produk.
type ProdukBarcodeBelumDipakaiValidator struct {
	barcode inservice.IBarcodeInService
}

var _ processor.IValidatorProcessor[dto.ProdukRequest] = (*ProdukBarcodeBelumDipakaiValidator)(nil)

func NewProdukBarcodeBelumDipakaiValidator(b inservice.IBarcodeInService) *ProdukBarcodeBelumDipakaiValidator {
	return &ProdukBarcodeBelumDipakaiValidator{barcode: b}
}

func (v *ProdukBarcodeBelumDipakaiValidator) Process(ctx context.Context, req dto.ProdukRequest) ([]processor.ValidationError, error) {
	kodes := make([]string, 0, len(req.Barcodes))
	for _, b := range req.Barcodes {
		kodes = append(kodes, b.Kode)
	}
	return cekBarcode(ctx, v.barcode, "barcodes", kodes)
}

// BarcodeBelumDipakaiValidator: barcode yang ditambahkan ke produk yang sudah ada.
type BarcodeBelumDipakaiValidator struct {
	barcode inservice.IBarcodeInService
}

var _ processor.IValidatorProcessor[dto.BarcodeAddRequest] = (*BarcodeBelumDipakaiValidator)(nil)

func NewBarcodeBelumDipakaiValidator(b inservice.IBarcodeInService) *BarcodeBelumDipakaiValidator {
	return &BarcodeBelumDipakaiValidator{barcode: b}
}

func (v *BarcodeBelumDipakaiValidator) Process(ctx context.Context, req dto.BarcodeAddRequest) ([]processor.ValidationError, error) {
	return cekBarcode(ctx, v.barcode, "kode", []string{req.Kode})
}
