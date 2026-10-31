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

type ProdukCreateExecutor struct {
	produk     inservice.IProdukInService
	actions    access.IActionChecker
	validators []processor.IValidatorProcessor[dto.ProdukRequest]
}

var _ processor.IExecutorProcessor[dto.ProdukRequest, dto.ProdukResponse] = (*ProdukCreateExecutor)(nil)

func NewProdukCreateExecutor(p inservice.IProdukInService, a access.IActionChecker, validators ...processor.IValidatorProcessor[dto.ProdukRequest]) *ProdukCreateExecutor {
	return &ProdukCreateExecutor{produk: p, actions: a, validators: validators}
}

func (e *ProdukCreateExecutor) Process(ctx context.Context, req dto.ProdukRequest) (dto.ProdukResponse, error) {
	req.Nama = strings.TrimSpace(req.Nama)
	if err := processor.Validate(ctx, req, e.validators...); err != nil {
		return dto.ProdukResponse{}, err
	}

	barcodes := make([]domain.BarcodeInput, 0, len(req.Barcodes))
	for _, b := range req.Barcodes {
		barcodes = append(barcodes, domain.BarcodeInput{Kode: b.Kode, IsiPerKemasan: b.IsiPerKemasan, Sumber: b.Sumber})
	}
	id, err := e.produk.Create(ctx, toProdukInput(req, true), barcodes)
	if err != nil {
		return dto.ProdukResponse{}, err
	}

	p, err := e.produk.GetByID(ctx, id)
	if err != nil {
		return dto.ProdukResponse{}, err
	}
	return mapper.ToProdukResponse(p, canSeeHPP(e.actions, req.RoleKode)), nil
}

func toProdukInput(req dto.ProdukRequest, aktif bool) domain.ProdukInput {
	return domain.ProdukInput{
		KategoriID:      req.KategoriID,
		Nama:            req.Nama,
		HargaNormal:     req.HargaNormal,
		BisaDitawar:     req.BisaDitawar,
		MarginMinPersen: req.MarginMinPersen,
		AturanPakai:     req.AturanPakai,
		Aktif:           aktif,
	}
}
