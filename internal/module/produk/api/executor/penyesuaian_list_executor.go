package executor

import (
	"context"

	"pos-jamu-backend/internal/common/access"
	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/produk/api/mapper"
	"pos-jamu-backend/internal/module/produk/tck/inservice"
)

type PenyesuaianListExecutor struct {
	penyesuaian inservice.IPenyesuaianStokInService
	produk      inservice.IProdukInService
	actions     access.IActionChecker
}

var _ processor.IExecutorProcessor[dto.PenyesuaianListRequest, []dto.PenyesuaianResponse] = (*PenyesuaianListExecutor)(nil)

func NewPenyesuaianListExecutor(ps inservice.IPenyesuaianStokInService, p inservice.IProdukInService, a access.IActionChecker) *PenyesuaianListExecutor {
	return &PenyesuaianListExecutor{penyesuaian: ps, produk: p, actions: a}
}

func (e *PenyesuaianListExecutor) Process(ctx context.Context, req dto.PenyesuaianListRequest) ([]dto.PenyesuaianResponse, error) {
	if _, err := e.produk.GetByID(ctx, req.ProdukID); err != nil {
		return nil, err
	}
	list, err := e.penyesuaian.ListByProduk(ctx, req.ProdukID, limitOrDefault(req.Limit))
	if err != nil {
		return nil, err
	}
	see := canSeeHPP(e.actions, req.RoleKode)
	out := make([]dto.PenyesuaianResponse, 0, len(list))
	for _, ps := range list {
		out = append(out, mapper.ToPenyesuaianResponse(ps, see))
	}
	return out, nil
}
