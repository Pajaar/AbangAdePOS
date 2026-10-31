package inservice

import (
	"context"
	"database/sql"

	"pos-jamu-backend/internal/module/auth/service/repository"
	"pos-jamu-backend/internal/module/auth/tck/domain"
	tck "pos-jamu-backend/internal/module/auth/tck/inservice"
	"pos-jamu-backend/internal/util/apperror"
)

type MenuInService struct {
	db *sql.DB
}

var _ tck.IMenuInService = (*MenuInService)(nil)

func NewMenuInService(db *sql.DB) *MenuInService {
	return &MenuInService{db: db}
}

func (s *MenuInService) ListByRole(ctx context.Context, roleKode string) ([]domain.Menu, error) {
	menus, err := repository.NewMenuRepository(s.db).ListByRole(ctx, roleKode)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return menus, nil
}
