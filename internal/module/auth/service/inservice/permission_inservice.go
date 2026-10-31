package inservice

import (
	"context"
	"database/sql"

	"pos-jamu-backend/internal/module/auth/service/repository"
	"pos-jamu-backend/internal/module/auth/tck/domain"
	tck "pos-jamu-backend/internal/module/auth/tck/inservice"
)

type PermissionInService struct {
	db *sql.DB
}

var _ tck.IPermissionInService = (*PermissionInService)(nil)

func NewPermissionInService(db *sql.DB) *PermissionInService {
	return &PermissionInService{db: db}
}

func (s *PermissionInService) ListAll(ctx context.Context) ([]domain.Permission, error) {
	return repository.NewPermissionRepository(s.db).ListAll(ctx)
}

func (s *PermissionInService) ListAllActions(ctx context.Context) ([]domain.RoleAction, error) {
	return repository.NewPermissionRepository(s.db).ListAllActions(ctx)
}
