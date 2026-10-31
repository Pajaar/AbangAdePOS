package executor

import (
	"context"

	"pos-jamu-backend/internal/common/access"
	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/produk/api/mapper"
	"pos-jamu-backend/internal/module/produk/tck/inservice"
)

type ProdukDetailExecutor struct {
	produk  inservice.IProdukInService
	actions access.IActionChecker
}

var _ processor.IExecutorProcessor[dto.IDRequest, dto.ProdukResponse] = (*ProdukDetailExecutor)(nil)

func NewProdukDetailExecutor(p inservice.IProdukInService, a access.IActionChecker) *ProdukDetailExecutor {
	return &ProdukDetailExecutor{produk: p, actions: a}
}

func (e *ProdukDetailExecutor) Process(ctx context.Context, req dto.IDRequest) (dto.ProdukResponse, error) {
	p, err := e.produk.GetByID(ctx, req.ID)
	if err != nil {
		return dto.ProdukResponse{}, err
	}
	return mapper.ToProdukResponse(p, canSeeHPP(e.actions, req.RoleKode)), nil
}
