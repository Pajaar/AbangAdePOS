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

type JamuUpdateExecutor struct {
	jamu       inservice.IJamuInService
	actions    access.IActionChecker
	validators []processor.IValidatorProcessor[dto.JamuRequest]
}

var _ processor.IExecutorProcessor[dto.JamuRequest, dto.JamuResponse] = (*JamuUpdateExecutor)(nil)

func NewJamuUpdateExecutor(j inservice.IJamuInService, a access.IActionChecker, validators ...processor.IValidatorProcessor[dto.JamuRequest]) *JamuUpdateExecutor {
	return &JamuUpdateExecutor{jamu: j, actions: a, validators: validators}
}

func (e *JamuUpdateExecutor) Process(ctx context.Context, req dto.JamuRequest) (dto.JamuResponse, error) {
	existing, err := e.jamu.GetByID(ctx, req.ID)
	if err != nil {
		return dto.JamuResponse{}, err
	}
	req.Nama = strings.TrimSpace(req.Nama)
	if err := processor.Validate(ctx, req, e.validators...); err != nil {
		return dto.JamuResponse{}, err
	}
	aktif := existing.Aktif
	if req.Aktif != nil {
		aktif = *req.Aktif
	}
	if err := e.jamu.Update(ctx, req.ID, domain.JamuInput{
		Nama: req.Nama, HargaJual: req.HargaJual, ProdukID: req.ProdukID, Aktif: aktif,
	}); err != nil {
		return dto.JamuResponse{}, err
	}
	j, err := e.jamu.GetByID(ctx, req.ID)
	if err != nil {
		return dto.JamuResponse{}, err
	}
	return mapper.ToJamuResponse(j, canSeeHPP(e.actions, req.RoleKode)), nil
}
