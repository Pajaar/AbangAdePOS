// Package executor berisi satu executor per use case modul auth (setara
// api/processor/executor di primecash). Executor hanya bergantung pada
// kontrak tck/inservice, bukan implementasinya.
package executor

import (
	"context"
	"time"

	"pos-jamu-backend/internal/auth"
	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/auth/api/mapper"
	"pos-jamu-backend/internal/module/auth/tck/domain"
	"pos-jamu-backend/internal/module/auth/tck/inservice"
	"pos-jamu-backend/internal/util/apperror"
)

// AccessTokenConfig dipakai executor yang menerbitkan access token (JWT).
type AccessTokenConfig struct {
	Secret string
	TTL    time.Duration
}

type LoginExecutor struct {
	userInService    inservice.IUserInService
	sessionInService inservice.ISessionInService
	accessToken      AccessTokenConfig
}

var _ processor.IExecutorProcessor[dto.LoginRequest, dto.LoginResult] = (*LoginExecutor)(nil)

func NewLoginExecutor(users inservice.IUserInService, sessions inservice.ISessionInService, at AccessTokenConfig) *LoginExecutor {
	return &LoginExecutor{userInService: users, sessionInService: sessions, accessToken: at}
}

func (e *LoginExecutor) Process(ctx context.Context, req dto.LoginRequest) (dto.LoginResult, error) {
	user, err := e.userInService.Authenticate(ctx, req.Username, req.Password)
	if err != nil {
		return dto.LoginResult{}, err
	}

	session, err := e.sessionInService.Start(ctx, user, domain.ClientInfo{
		UserAgent: req.UserAgent,
		IPAddress: req.IPAddress,
	})
	if err != nil {
		return dto.LoginResult{}, err
	}

	accessToken, expiresAt, err := issueAccessToken(user.ID, user.RoleKode, e.accessToken)
	if err != nil {
		return dto.LoginResult{}, err
	}

	return dto.LoginResult{
		Response: dto.LoginResponse{
			AccessToken: accessToken,
			ExpiresAt:   expiresAt,
			User:        mapper.ToUserInfo(user),
		},
		RefreshToken:     session.RefreshToken,
		RefreshExpiresAt: session.RefreshExpiresAt,
	}, nil
}

func issueAccessToken(userID uint64, roleKode string, cfg AccessTokenConfig) (string, time.Time, error) {
	token, err := auth.GenerateAccessToken(userID, roleKode, cfg.Secret, cfg.TTL)
	if err != nil {
		return "", time.Time{}, apperror.Internal(err)
	}
	return token, time.Now().Add(cfg.TTL), nil
}
