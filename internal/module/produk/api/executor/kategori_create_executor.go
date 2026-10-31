package executor

import (
	"context"
	"strings"

	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/produk/api/mapper"
	"pos-jamu-backend/internal/module/produk/tck/inservice"
)

type KategoriCreateExecutor struct {
	kategori   inservice.IKategoriInService
	validators []processor.IValidatorProcessor[dto.KategoriRequest]
}

var _ processor.IExecutorProcessor[dto.KategoriRequest, dto.KategoriResponse] = (*KategoriCreateExecutor)(nil)

func NewKategoriCreateExecutor(k inservice.IKategoriInService, validators ...processor.IValidatorProcessor[dto.KategoriRequest]) *KategoriCreateExecutor {
	return &KategoriCreateExecutor{kategori: k, validators: validators}
}

func (e *KategoriCreateExecutor) Process(ctx context.Context, req dto.KategoriRequest) (dto.KategoriResponse, error) {
	req.Nama = strings.TrimSpace(req.Nama)
	if err := processor.Validate(ctx, req, e.validators...); err != nil {
		return dto.KategoriResponse{}, err
	}
	id, err := e.kategori.Create(ctx, req.Nama)
	if err != nil {
		return dto.KategoriResponse{}, err
	}
	k, err := e.kategori.GetByID(ctx, id)
	if err != nil {
		return dto.KategoriResponse{}, err
	}
	return mapper.ToKategoriResponse(k), nil
}
