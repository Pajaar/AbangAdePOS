package inservice

import (
	"context"
	"database/sql"
	"errors"

	"pos-jamu-backend/internal/module/produk/service/repository"
	"pos-jamu-backend/internal/module/produk/tck/domain"
	tck "pos-jamu-backend/internal/module/produk/tck/inservice"
	"pos-jamu-backend/internal/util/apperror"
)

type KategoriInService struct {
	db *sql.DB
}

var _ tck.IKategoriInService = (*KategoriInService)(nil)

func NewKategoriInService(db *sql.DB) *KategoriInService {
	return &KategoriInService{db: db}
}

func (s *KategoriInService) List(ctx context.Context, aktif *bool) ([]domain.Kategori, error) {
	list, err := repository.NewKategoriRepository(s.db).List(ctx, aktif)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return list, nil
}

func (s *KategoriInService) GetByID(ctx context.Context, id uint64) (domain.Kategori, error) {
	k, err := repository.NewKategoriRepository(s.db).GetByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Kategori{}, apperror.NotFound("Kategori tidak ditemukan")
	}
	if err != nil {
		return domain.Kategori{}, apperror.Internal(err)
	}
	return k, nil
}

func (s *KategoriInService) FindByNama(ctx context.Context, nama string) (domain.Kategori, bool, error) {
	k, err := repository.NewKategoriRepository(s.db).GetByNama(ctx, nama)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Kategori{}, false, nil
	}
	if err != nil {
		return domain.Kategori{}, false, apperror.Internal(err)
	}
	return k, true, nil
}

func (s *KategoriInService) Create(ctx context.Context, nama string) (uint64, error) {
	id, err := repository.NewKategoriRepository(s.db).Create(ctx, nama)
	if mysqlErrNo(err) == mysqlErrDuplicate {
		return 0, apperror.Conflict("Nama kategori sudah dipakai")
	}
	if err != nil {
		return 0, apperror.Internal(err)
	}
	return id, nil
}

func (s *KategoriInService) Update(ctx context.Context, id uint64, nama string, aktif bool) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return err
	}
	err := repository.NewKategoriRepository(s.db).Update(ctx, id, nama, aktif)
	if mysqlErrNo(err) == mysqlErrDuplicate {
		return apperror.Conflict("Nama kategori sudah dipakai")
	}
	if err != nil {
		return apperror.Internal(err)
	}
	return nil
}
