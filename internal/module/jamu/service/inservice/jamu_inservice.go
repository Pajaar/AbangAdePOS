// Package inservice berisi implementasi kontrak tck/inservice modul menu jamu.
package inservice

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"

	"pos-jamu-backend/internal/module/jamu/service/repository"
	"pos-jamu-backend/internal/module/jamu/tck/domain"
	tck "pos-jamu-backend/internal/module/jamu/tck/inservice"
	"pos-jamu-backend/internal/util/apperror"
)

const mysqlErrForeignKey = 1452

// mapWriteErr: produk yang ditunjuk tidak ada (FK) -> 400.
func mapWriteErr(err error) error {
	var me *mysql.MySQLError
	if errors.As(err, &me) && me.Number == mysqlErrForeignKey {
		return apperror.BadRequest("Produk tidak ditemukan")
	}
	return apperror.Internal(err)
}

type JamuInService struct {
	db *sql.DB
}

var _ tck.IJamuInService = (*JamuInService)(nil)

func NewJamuInService(db *sql.DB) *JamuInService {
	return &JamuInService{db: db}
}

func (s *JamuInService) List(ctx context.Context, aktif *bool) ([]domain.Jamu, error) {
	list, err := repository.NewJamuRepository(s.db).List(ctx, aktif)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return list, nil
}

func (s *JamuInService) GetByID(ctx context.Context, id uint64) (domain.Jamu, error) {
	j, err := repository.NewJamuRepository(s.db).GetByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Jamu{}, apperror.NotFound("Jamu tidak ditemukan")
	}
	if err != nil {
		return domain.Jamu{}, apperror.Internal(err)
	}
	return j, nil
}

func (s *JamuInService) Create(ctx context.Context, in domain.JamuInput) (uint64, error) {
	id, err := repository.NewJamuRepository(s.db).Create(ctx, in)
	if err != nil {
		return 0, mapWriteErr(err)
	}
	return id, nil
}

func (s *JamuInService) Update(ctx context.Context, id uint64, in domain.JamuInput) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return err
	}
	if err := repository.NewJamuRepository(s.db).Update(ctx, id, in); err != nil {
		return mapWriteErr(err)
	}
	return nil
}
