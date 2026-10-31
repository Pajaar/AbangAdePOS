package executor

import (
	"context"

	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/produk/api/mapper"
	"pos-jamu-backend/internal/module/produk/tck/inservice"
)

type KategoriListExecutor struct {
	kategori inservice.IKategoriInService
}

var _ processor.IExecutorProcessor[dto.KategoriListRequest, []dto.KategoriResponse] = (*KategoriListExecutor)(nil)

func NewKategoriListExecutor(k inservice.IKategoriInService) *KategoriListExecutor {
	return &KategoriListExecutor{kategori: k}
}

func (e *KategoriListExecutor) Process(ctx context.Context, req dto.KategoriListRequest) ([]dto.KategoriResponse, error) {
	list, err := e.kategori.List(ctx, req.Aktif)
	if err != nil {
		return nil, err
	}
	out := make([]dto.KategoriResponse, 0, len(list))
	for _, k := range list {
		out = append(out, mapper.ToKategoriResponse(k))
	}
	return out, nil
}
