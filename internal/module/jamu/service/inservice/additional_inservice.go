package inservice

import (
	"context"
	"database/sql"
	"errors"

	"pos-jamu-backend/internal/module/jamu/service/repository"
	"pos-jamu-backend/internal/module/jamu/tck/domain"
	tck "pos-jamu-backend/internal/module/jamu/tck/inservice"
	"pos-jamu-backend/internal/util/apperror"
)

type AdditionalItemInService struct {
	db *sql.DB
}

var _ tck.IAdditionalItemInService = (*AdditionalItemInService)(nil)

func NewAdditionalItemInService(db *sql.DB) *AdditionalItemInService {
	return &AdditionalItemInService{db: db}
}

func (s *AdditionalItemInService) List(ctx context.Context, aktif *bool) ([]domain.AdditionalItem, error) {
	list, err := repository.NewAdditionalRepository(s.db).List(ctx, aktif)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return list, nil
}

func (s *AdditionalItemInService) GetByID(ctx context.Context, id uint64) (domain.AdditionalItem, error) {
	a, err := repository.NewAdditionalRepository(s.db).GetByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.AdditionalItem{}, apperror.NotFound("Item tambahan tidak ditemukan")
	}
	if err != nil {
		return domain.AdditionalItem{}, apperror.Internal(err)
	}
	return a, nil
}

func (s *AdditionalItemInService) Create(ctx context.Context, in domain.AdditionalInput) (uint64, error) {
	id, err := repository.NewAdditionalRepository(s.db).Create(ctx, in)
	if err != nil {
		return 0, mapWriteErr(err)
	}
	return id, nil
}

func (s *AdditionalItemInService) Update(ctx context.Context, id uint64, in domain.AdditionalInput) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return err
	}
	if err := repository.NewAdditionalRepository(s.db).Update(ctx, id, in); err != nil {
		return mapWriteErr(err)
	}
	return nil
}
