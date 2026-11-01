package executor

import (
	"context"
	"strings"

	"pos-jamu-backend/internal/common/access"
	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/jamu/api/mapper"
	"pos-jamu-backend/internal/module/jamu/tck/domain"
	"pos-jamu-backend/internal/module/jamu/tck/inservice"
)

type AdditionalUpdateExecutor struct {
	additional inservice.IAdditionalItemInService
	actions    access.IActionChecker
	validators []processor.IValidatorProcessor[dto.AdditionalRequest]
}

var _ processor.IExecutorProcessor[dto.AdditionalRequest, dto.AdditionalResponse] = (*AdditionalUpdateExecutor)(nil)

func NewAdditionalUpdateExecutor(s inservice.IAdditionalItemInService, a access.IActionChecker, validators ...processor.IValidatorProcessor[dto.AdditionalRequest]) *AdditionalUpdateExecutor {
	return &AdditionalUpdateExecutor{additional: s, actions: a, validators: validators}
}

func (e *AdditionalUpdateExecutor) Process(ctx context.Context, req dto.AdditionalRequest) (dto.AdditionalResponse, error) {
	existing, err := e.additional.GetByID(ctx, req.ID)
	if err != nil {
		return dto.AdditionalResponse{}, err
	}
	req.Nama = strings.TrimSpace(req.Nama)
	if err := processor.Validate(ctx, req, e.validators...); err != nil {
		return dto.AdditionalResponse{}, err
	}
	aktif := existing.Aktif
	if req.Aktif != nil {
		aktif = *req.Aktif
	}
	// produk_id kosong = tambahan tidak dilacak stoknya.
	if err := e.additional.Update(ctx, req.ID, domain.AdditionalInput{
		Nama: req.Nama, HargaJual: req.HargaJual, ProdukID: req.ProdukID, Aktif: aktif,
	}); err != nil {
		return dto.AdditionalResponse{}, err
	}
	a, err := e.additional.GetByID(ctx, req.ID)
	if err != nil {
		return dto.AdditionalResponse{}, err
	}
	return mapper.ToAdditionalResponse(a, canSeeHPP(e.actions, req.RoleKode)), nil
}
