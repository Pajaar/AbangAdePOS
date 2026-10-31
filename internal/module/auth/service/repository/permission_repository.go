package repository

import (
	"context"

	"pos-jamu-backend/internal/module/auth/tck/domain"
	sqlcgen "pos-jamu-backend/internal/repository/sqlc"
)

type IPermissionRepository interface {
	ListAll(ctx context.Context) ([]domain.Permission, error)
	ListAllActions(ctx context.Context) ([]domain.RoleAction, error)
}

type permissionRepository struct {
	q *sqlcgen.Queries
}

func NewPermissionRepository(db sqlcgen.DBTX) IPermissionRepository {
	return &permissionRepository{q: sqlcgen.New(db)}
}

func (r *permissionRepository) ListAll(ctx context.Context) ([]domain.Permission, error) {
	rows, err := r.q.ListAllPermissions(ctx)
	if err != nil {
		return nil, err
	}
	perms := make([]domain.Permission, 0, len(rows))
	for _, row := range rows {
		perms = append(perms, domain.Permission{
			RoleKode: row.RoleKode,
			Method:   string(row.Method),
			Path:     row.Path,
		})
	}
	return perms, nil
}

func (r *permissionRepository) ListAllActions(ctx context.Context) ([]domain.RoleAction, error) {
	rows, err := r.q.ListAllRoleActions(ctx)
	if err != nil {
		return nil, err
	}
	actions := make([]domain.RoleAction, 0, len(rows))
	for _, row := range rows {
		actions = append(actions, domain.RoleAction{RoleKode: row.RoleKode, AksiKode: row.AksiKode})
	}
	return actions, nil
}
