package executor

import (
	"context"

	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/auth/api/mapper"
	"pos-jamu-backend/internal/module/auth/tck/inservice"
)

type MenusExecutor struct {
	menuInService inservice.IMenuInService
}

var _ processor.IExecutorProcessor[dto.MenusRequest, dto.MenusResponse] = (*MenusExecutor)(nil)

func NewMenusExecutor(menus inservice.IMenuInService) *MenusExecutor {
	return &MenusExecutor{menuInService: menus}
}

func (e *MenusExecutor) Process(ctx context.Context, req dto.MenusRequest) (dto.MenusResponse, error) {
	menus, err := e.menuInService.ListByRole(ctx, req.RoleKode)
	if err != nil {
		return dto.MenusResponse{}, err
	}
	return mapper.ToMenusResponse(menus), nil
}
