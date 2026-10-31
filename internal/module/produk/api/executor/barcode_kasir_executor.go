package executor

import (
	"context"

	"pos-jamu-backend/internal/common/access"
	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/produk/api/mapper"
	"pos-jamu-backend/internal/module/produk/tck/inservice"
)

// BarcodeKasirExecutor: scan di kasir, hanya database lokal (tanpa API publik).
type BarcodeKasirExecutor struct {
	barcode inservice.IBarcodeInService
	actions access.IActionChecker
}

var _ processor.IExecutorProcessor[dto.BarcodeKasirRequest, dto.BarcodeKasirResponse] = (*BarcodeKasirExecutor)(nil)

func NewBarcodeKasirExecutor(b inservice.IBarcodeInService, a access.IActionChecker) *BarcodeKasirExecutor {
	return &BarcodeKasirExecutor{barcode: b, actions: a}
}

func (e *BarcodeKasirExecutor) Process(ctx context.Context, req dto.BarcodeKasirRequest) (dto.BarcodeKasirResponse, error) {
	p, b, err := e.barcode.FindForKasir(ctx, req.Kode)
	if err != nil {
		return dto.BarcodeKasirResponse{}, err
	}
	p.Barcodes = nil // kasir cukup barcode yang di-scan
	return dto.BarcodeKasirResponse{
		Produk:  mapper.ToProdukResponse(p, canSeeHPP(e.actions, req.RoleKode)),
		Barcode: mapper.ToBarcodeResponse(b),
	}, nil
}
