package executor

import (
	"context"

	"pos-jamu-backend/internal/common/access"
	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/produk/api/mapper"
	"pos-jamu-backend/internal/module/produk/tck/inservice"
)

// BarcodeLookupExecutor: dipakai di form tambah produk / pembelian.
type BarcodeLookupExecutor struct {
	barcode inservice.IBarcodeInService
	actions access.IActionChecker
}

var _ processor.IExecutorProcessor[dto.BarcodeLookupRequest, dto.BarcodeLookupResponse] = (*BarcodeLookupExecutor)(nil)

func NewBarcodeLookupExecutor(b inservice.IBarcodeInService, a access.IActionChecker) *BarcodeLookupExecutor {
	return &BarcodeLookupExecutor{barcode: b, actions: a}
}

func (e *BarcodeLookupExecutor) Process(ctx context.Context, req dto.BarcodeLookupRequest) (dto.BarcodeLookupResponse, error) {
	r, err := e.barcode.Lookup(ctx, req.Kode)
	if err != nil {
		return dto.BarcodeLookupResponse{}, err
	}
	return mapper.ToLookupResponse(r, canSeeHPP(e.actions, req.RoleKode)), nil
}
