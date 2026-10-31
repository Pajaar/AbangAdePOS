package executor

import (
	"context"

	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/auth/tck/domain"
	"pos-jamu-backend/internal/module/auth/tck/inservice"
)

type RefreshExecutor struct {
	sessionInService inservice.ISessionInService
	accessToken      AccessTokenConfig
}

var _ processor.IExecutorProcessor[dto.RefreshRequest, dto.RefreshResult] = (*RefreshExecutor)(nil)

func NewRefreshExecutor(sessions inservice.ISessionInService, at AccessTokenConfig) *RefreshExecutor {
	return &RefreshExecutor{sessionInService: sessions, accessToken: at}
}

func (e *RefreshExecutor) Process(ctx context.Context, req dto.RefreshRequest) (dto.RefreshResult, error) {
	session, err := e.sessionInService.Rotate(ctx, req.RefreshToken, domain.ClientInfo{
		UserAgent: req.UserAgent,
		IPAddress: req.IPAddress,
	})
	if err != nil {
		return dto.RefreshResult{}, err
	}

	accessToken, expiresAt, err := issueAccessToken(session.UserID, session.RoleKode, e.accessToken)
	if err != nil {
		return dto.RefreshResult{}, err
	}

	return dto.RefreshResult{
		Response: dto.RefreshResponse{
			AccessToken: accessToken,
			ExpiresAt:   expiresAt,
		},
		RefreshToken:     session.RefreshToken,
		RefreshExpiresAt: session.RefreshExpiresAt,
	}, nil
}
