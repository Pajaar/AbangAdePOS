package executor

import (
	"context"
	"errors"
	"net/http"

	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/auth/api/mapper"
	"pos-jamu-backend/internal/module/auth/tck/inservice"
	"pos-jamu-backend/internal/util/apperror"
)

type MeExecutor struct {
	userInService inservice.IUserInService
}

var _ processor.IExecutorProcessor[dto.MeRequest, dto.MeResponse] = (*MeExecutor)(nil)

func NewMeExecutor(users inservice.IUserInService) *MeExecutor {
	return &MeExecutor{userInService: users}
}

func (e *MeExecutor) Process(ctx context.Context, req dto.MeRequest) (dto.MeResponse, error) {
	user, err := e.userInService.GetByID(ctx, req.UserID)
	if err != nil {
		// User di token sudah tidak ada = sesi tidak berlaku (401), bukan 404.
		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.Status == http.StatusNotFound {
			return dto.MeResponse{}, apperror.Unauthorized("User tidak ditemukan")
		}
		return dto.MeResponse{}, err
	}
	return dto.MeResponse{UserInfo: mapper.ToUserInfo(user)}, nil
}
