package repository

import (
	"context"

	sqlcgen "pos-jamu-backend/internal/repository/sqlc"
)

type IRoleRepository interface {
	// GetIDByKode mengembalikan sql.ErrNoRows bila role tidak ada.
	GetIDByKode(ctx context.Context, kode string) (uint64, error)
}

type roleRepository struct {
	q *sqlcgen.Queries
}

func NewRoleRepository(db sqlcgen.DBTX) IRoleRepository {
	return &roleRepository{q: sqlcgen.New(db)}
}

func (r *roleRepository) GetIDByKode(ctx context.Context, kode string) (uint64, error) {
	role, err := r.q.GetRoleByKode(ctx, kode)
	return role.ID, err
}
