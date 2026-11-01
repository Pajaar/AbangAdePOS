package executor

import (
	"context"

	"pos-jamu-backend/internal/common/access"
	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/jamu/api/mapper"
	"pos-jamu-backend/internal/module/jamu/tck/inservice"
)

type AdditionalListExecutor struct {
	additional inservice.IAdditionalItemInService
	actions    access.IActionChecker
}

var _ processor.IExecutorProcessor[dto.AdditionalListRequest, []dto.AdditionalResponse] = (*AdditionalListExecutor)(nil)

func NewAdditionalListExecutor(s inservice.IAdditionalItemInService, a access.IActionChecker) *AdditionalListExecutor {
	return &AdditionalListExecutor{additional: s, actions: a}
}

func (e *AdditionalListExecutor) Process(ctx context.Context, req dto.AdditionalListRequest) ([]dto.AdditionalResponse, error) {
	list, err := e.additional.List(ctx, req.Aktif)
	if err != nil {
		return nil, err
	}
	see := canSeeHPP(e.actions, req.RoleKode)
	out := make([]dto.AdditionalResponse, 0, len(list))
	for _, a := range list {
		out = append(out, mapper.ToAdditionalResponse(a, see))
	}
	return out, nil
}
