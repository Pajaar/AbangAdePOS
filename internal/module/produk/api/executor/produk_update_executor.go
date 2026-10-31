package executor

import (
	"context"
	"strings"

	"pos-jamu-backend/internal/common/access"
	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/produk/api/mapper"
	"pos-jamu-backend/internal/module/produk/tck/inservice"
)

type ProdukUpdateExecutor struct {
	produk     inservice.IProdukInService
	actions    access.IActionChecker
	validators []processor.IValidatorProcessor[dto.ProdukRequest]
}

var _ processor.IExecutorProcessor[dto.ProdukRequest, dto.ProdukResponse] = (*ProdukUpdateExecutor)(nil)

func NewProdukUpdateExecutor(p inservice.IProdukInService, a access.IActionChecker, validators ...processor.IValidatorProcessor[dto.ProdukRequest]) *ProdukUpdateExecutor {
	return &ProdukUpdateExecutor{produk: p, actions: a, validators: validators}
}

func (e *ProdukUpdateExecutor) Process(ctx context.Context, req dto.ProdukRequest) (dto.ProdukResponse, error) {
	existing, err := e.produk.GetByID(ctx, req.ID)
	if err != nil {
		return dto.ProdukResponse{}, err
	}
	see := canSeeHPP(e.actions, req.RoleKode)

	// Role tanpa lihat_hpp tidak pernah menerima margin, jadi margin yang
	// tidak dikirim berarti tidak diubah, bukan dihapus.
	if !see && req.MarginMinPersen == nil {
		req.MarginMinPersen = existing.MarginMinPersen
	}
	aktif := existing.Aktif
	if req.Aktif != nil {
		aktif = *req.Aktif
	}
	req.Nama = strings.TrimSpace(req.Nama)

	if err := processor.Validate(ctx, req, e.validators...); err != nil {
		return dto.ProdukResponse{}, err
	}
	if err := e.produk.Update(ctx, req.ID, toProdukInput(req, aktif)); err != nil {
		return dto.ProdukResponse{}, err
	}

	p, err := e.produk.GetByID(ctx, req.ID)
	if err != nil {
		return dto.ProdukResponse{}, err
	}
	return mapper.ToProdukResponse(p, see), nil
}
