package executor

import (
	"context"

	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/auth/tck/inservice"
)

type LogoutExecutor struct {
	sessionInService inservice.ISessionInService
}

var _ processor.IExecutorProcessor[dto.LogoutRequest, processor.Empty] = (*LogoutExecutor)(nil)

func NewLogoutExecutor(sessions inservice.ISessionInService) *LogoutExecutor {
	return &LogoutExecutor{sessionInService: sessions}
}

func (e *LogoutExecutor) Process(ctx context.Context, req dto.LogoutRequest) (processor.Empty, error) {
	return processor.Empty{}, e.sessionInService.End(ctx, req.RefreshToken)
}
