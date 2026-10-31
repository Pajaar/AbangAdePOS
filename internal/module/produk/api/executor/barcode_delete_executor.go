package executor

import (
	"context"

	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/produk/tck/inservice"
)

type BarcodeDeleteExecutor struct {
	barcode inservice.IBarcodeInService
}

var _ processor.IExecutorProcessor[dto.IDRequest, processor.Empty] = (*BarcodeDeleteExecutor)(nil)

func NewBarcodeDeleteExecutor(b inservice.IBarcodeInService) *BarcodeDeleteExecutor {
	return &BarcodeDeleteExecutor{barcode: b}
}

func (e *BarcodeDeleteExecutor) Process(ctx context.Context, req dto.IDRequest) (processor.Empty, error) {
	return processor.Empty{}, e.barcode.Delete(ctx, req.ID)
}
