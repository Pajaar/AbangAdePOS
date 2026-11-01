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

type AdditionalCreateExecutor struct {
	additional inservice.IAdditionalItemInService
	actions    access.IActionChecker
	validators []processor.IValidatorProcessor[dto.AdditionalRequest]
}

var _ processor.IExecutorProcessor[dto.AdditionalRequest, dto.AdditionalResponse] = (*AdditionalCreateExecutor)(nil)

func NewAdditionalCreateExecutor(s inservice.IAdditionalItemInService, a access.IActionChecker, validators ...processor.IValidatorProcessor[dto.AdditionalRequest]) *AdditionalCreateExecutor {
	return &AdditionalCreateExecutor{additional: s, actions: a, validators: validators}
}

func (e *AdditionalCreateExecutor) Process(ctx context.Context, req dto.AdditionalRequest) (dto.AdditionalResponse, error) {
	req.Nama = strings.TrimSpace(req.Nama)
	if err := processor.Validate(ctx, req, e.validators...); err != nil {
		return dto.AdditionalResponse{}, err
	}
	id, err := e.additional.Create(ctx, domain.AdditionalInput{Nama: req.Nama, HargaJual: req.HargaJual, ProdukID: req.ProdukID})
	if err != nil {
		return dto.AdditionalResponse{}, err
	}
	a, err := e.additional.GetByID(ctx, id)
	if err != nil {
		return dto.AdditionalResponse{}, err
	}
	return mapper.ToAdditionalResponse(a, canSeeHPP(e.actions, req.RoleKode)), nil
}
