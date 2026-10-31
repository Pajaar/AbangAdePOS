package repository

import (
	"context"

	"pos-jamu-backend/internal/module/auth/tck/domain"
	sqlcgen "pos-jamu-backend/internal/repository/sqlc"
)

type IMenuRepository interface {
	// ListByRole mengembalikan menu tipe 'menu' dan 'aksi' milik role.
	ListByRole(ctx context.Context, roleKode string) ([]domain.Menu, error)
}

type menuRepository struct {
	q *sqlcgen.Queries
}

func NewMenuRepository(db sqlcgen.DBTX) IMenuRepository {
	return &menuRepository{q: sqlcgen.New(db)}
}

func (r *menuRepository) ListByRole(ctx context.Context, roleKode string) ([]domain.Menu, error) {
	rows, err := r.q.ListMenusByRole(ctx, roleKode)
	if err != nil {
		return nil, err
	}
	menus := make([]domain.Menu, 0, len(rows))
	for _, row := range rows {
		m := domain.Menu{
			ID:     row.ID,
			Kode:   row.Kode,
			Nama:   row.Nama,
			Tipe:   domain.MenuTipe(row.Tipe),
			Urutan: row.Urutan,
		}
		if row.ParentID.Valid {
			p := uint64(row.ParentID.Int64)
			m.ParentID = &p
		}
		if row.Path.Valid {
			m.Path = &row.Path.String
		}
		if row.Icon.Valid {
			m.Icon = &row.Icon.String
		}
		menus = append(menus, m)
	}
	return menus, nil
}
