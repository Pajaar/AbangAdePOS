package executor

import (
	"context"
	"strings"

	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/produk/api/mapper"
	"pos-jamu-backend/internal/module/produk/tck/inservice"
)

type KategoriUpdateExecutor struct {
	kategori   inservice.IKategoriInService
	validators []processor.IValidatorProcessor[dto.KategoriRequest]
}

var _ processor.IExecutorProcessor[dto.KategoriRequest, dto.KategoriResponse] = (*KategoriUpdateExecutor)(nil)

func NewKategoriUpdateExecutor(k inservice.IKategoriInService, validators ...processor.IValidatorProcessor[dto.KategoriRequest]) *KategoriUpdateExecutor {
	return &KategoriUpdateExecutor{kategori: k, validators: validators}
}

func (e *KategoriUpdateExecutor) Process(ctx context.Context, req dto.KategoriRequest) (dto.KategoriResponse, error) {
	existing, err := e.kategori.GetByID(ctx, req.ID)
	if err != nil {
		return dto.KategoriResponse{}, err
	}
	req.Nama = strings.TrimSpace(req.Nama)
	if err := processor.Validate(ctx, req, e.validators...); err != nil {
		return dto.KategoriResponse{}, err
	}
	aktif := existing.Aktif // aktif tidak dikirim = tidak berubah
	if req.Aktif != nil {
		aktif = *req.Aktif
	}
	if err := e.kategori.Update(ctx, req.ID, req.Nama, aktif); err != nil {
		return dto.KategoriResponse{}, err
	}
	k, err := e.kategori.GetByID(ctx, req.ID)
	if err != nil {
		return dto.KategoriResponse{}, err
	}
	return mapper.ToKategoriResponse(k), nil
}
