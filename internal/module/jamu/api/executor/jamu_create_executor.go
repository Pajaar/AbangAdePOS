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

type JamuCreateExecutor struct {
	jamu       inservice.IJamuInService
	actions    access.IActionChecker
	validators []processor.IValidatorProcessor[dto.JamuRequest]
}

var _ processor.IExecutorProcessor[dto.JamuRequest, dto.JamuResponse] = (*JamuCreateExecutor)(nil)

func NewJamuCreateExecutor(j inservice.IJamuInService, a access.IActionChecker, validators ...processor.IValidatorProcessor[dto.JamuRequest]) *JamuCreateExecutor {
	return &JamuCreateExecutor{jamu: j, actions: a, validators: validators}
}

func (e *JamuCreateExecutor) Process(ctx context.Context, req dto.JamuRequest) (dto.JamuResponse, error) {
	req.Nama = strings.TrimSpace(req.Nama)
	if err := processor.Validate(ctx, req, e.validators...); err != nil {
		return dto.JamuResponse{}, err
	}
	id, err := e.jamu.Create(ctx, domain.JamuInput{Nama: req.Nama, HargaJual: req.HargaJual, ProdukID: req.ProdukID})
	if err != nil {
		return dto.JamuResponse{}, err
	}
	j, err := e.jamu.GetByID(ctx, id)
	if err != nil {
		return dto.JamuResponse{}, err
	}
	return mapper.ToJamuResponse(j, canSeeHPP(e.actions, req.RoleKode)), nil
}
