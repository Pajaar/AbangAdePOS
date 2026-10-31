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

type PenyesuaianCreateExecutor struct {
	penyesuaian inservice.IPenyesuaianStokInService
	actions     access.IActionChecker
}

var _ processor.IExecutorProcessor[dto.PenyesuaianRequest, dto.PenyesuaianResponse] = (*PenyesuaianCreateExecutor)(nil)

func NewPenyesuaianCreateExecutor(p inservice.IPenyesuaianStokInService, a access.IActionChecker) *PenyesuaianCreateExecutor {
	return &PenyesuaianCreateExecutor{penyesuaian: p, actions: a}
}

func (e *PenyesuaianCreateExecutor) Process(ctx context.Context, req dto.PenyesuaianRequest) (dto.PenyesuaianResponse, error) {
	ps, err := e.penyesuaian.Create(ctx, domain.PenyesuaianInput{
		ProdukID:   req.ProdukID,
		UserID:     req.UserID,
		Jenis:      domain.JenisPenyesuaian(req.Jenis),
		Qty:        req.Qty,
		StokFisik:  req.StokFisik,
		HargaModal: req.HargaModal,
		Keterangan: strings.TrimSpace(req.Keterangan),
	})
	if err != nil {
		return dto.PenyesuaianResponse{}, err
	}
	return mapper.ToPenyesuaianResponse(ps, canSeeHPP(e.actions, req.RoleKode)), nil
}
