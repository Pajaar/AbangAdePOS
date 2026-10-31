package executor

import (
	"context"
	"strings"

	"pos-jamu-backend/internal/common/access"
	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/produk/api/mapper"
	"pos-jamu-backend/internal/module/produk/tck/domain"
	"pos-jamu-backend/internal/module/produk/tck/inservice"
)

type ProdukListExecutor struct {
	produk  inservice.IProdukInService
	actions access.IActionChecker
}

var _ processor.IExecutorProcessor[dto.ProdukListRequest, dto.PagedResponse[dto.ProdukResponse]] = (*ProdukListExecutor)(nil)

func NewProdukListExecutor(p inservice.IProdukInService, a access.IActionChecker) *ProdukListExecutor {
	return &ProdukListExecutor{produk: p, actions: a}
}

func (e *ProdukListExecutor) Process(ctx context.Context, req dto.ProdukListRequest) (dto.PagedResponse[dto.ProdukResponse], error) {
	f := domain.ProdukFilter{
		KategoriID: req.KategoriID,
		Aktif:      req.Aktif,
		Limit:      limitOrDefault(req.Limit),
		Offset:     req.Offset,
	}
	if q := strings.TrimSpace(req.Q); q != "" {
		f.Q = &q
	}
	if req.Status != "" {
		st := domain.Status(req.Status)
		f.Status = &st
	}

	items, total, err := e.produk.List(ctx, f)
	if err != nil {
		return dto.PagedResponse[dto.ProdukResponse]{}, err
	}
	see := canSeeHPP(e.actions, req.RoleKode)
	out := make([]dto.ProdukResponse, 0, len(items))
	for _, p := range items {
		out = append(out, mapper.ToProdukResponse(p, see))
	}
	return dto.PagedResponse[dto.ProdukResponse]{Items: out, Total: total, Limit: f.Limit, Offset: f.Offset}, nil
}
