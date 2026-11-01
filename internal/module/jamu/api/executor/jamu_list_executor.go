package executor

import (
	"context"

	"pos-jamu-backend/internal/common/access"
	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/jamu/api/mapper"
	"pos-jamu-backend/internal/module/jamu/tck/inservice"
)

type JamuListExecutor struct {
	jamu    inservice.IJamuInService
	actions access.IActionChecker
}

var _ processor.IExecutorProcessor[dto.JamuListRequest, []dto.JamuResponse] = (*JamuListExecutor)(nil)

func NewJamuListExecutor(j inservice.IJamuInService, a access.IActionChecker) *JamuListExecutor {
	return &JamuListExecutor{jamu: j, actions: a}
}

func (e *JamuListExecutor) Process(ctx context.Context, req dto.JamuListRequest) ([]dto.JamuResponse, error) {
	list, err := e.jamu.List(ctx, req.Aktif)
	if err != nil {
		return nil, err
	}
	see := canSeeHPP(e.actions, req.RoleKode)
	out := make([]dto.JamuResponse, 0, len(list))
	for _, j := range list {
		out = append(out, mapper.ToJamuResponse(j, see))
	}
	return out, nil
}
