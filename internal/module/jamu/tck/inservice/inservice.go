// Package inservice berisi kontrak service modul menu jamu.
package inservice

import (
	"context"

	"pos-jamu-backend/internal/module/jamu/tck/domain"
)

type IJamuInService interface {
	List(ctx context.Context, aktif *bool) ([]domain.Jamu, error)
	// GetByID mengembalikan apperror NotFound bila tidak ada.
	GetByID(ctx context.Context, id uint64) (domain.Jamu, error)
	Create(ctx context.Context, in domain.JamuInput) (uint64, error)
	Update(ctx context.Context, id uint64, in domain.JamuInput) error
}

type IAdditionalItemInService interface {
	List(ctx context.Context, aktif *bool) ([]domain.AdditionalItem, error)
	GetByID(ctx context.Context, id uint64) (domain.AdditionalItem, error)
	Create(ctx context.Context, in domain.AdditionalInput) (uint64, error)
	Update(ctx context.Context, id uint64, in domain.AdditionalInput) error
}
