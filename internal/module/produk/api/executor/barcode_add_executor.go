package executor

import (
	"context"

	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/produk/api/mapper"
	"pos-jamu-backend/internal/module/produk/tck/domain"
	"pos-jamu-backend/internal/module/produk/tck/inservice"
)

type BarcodeAddExecutor struct {
	barcode    inservice.IBarcodeInService
	validators []processor.IValidatorProcessor[dto.BarcodeAddRequest]
}

var _ processor.IExecutorProcessor[dto.BarcodeAddRequest, dto.BarcodeResponse] = (*BarcodeAddExecutor)(nil)

func NewBarcodeAddExecutor(b inservice.IBarcodeInService, validators ...processor.IValidatorProcessor[dto.BarcodeAddRequest]) *BarcodeAddExecutor {
	return &BarcodeAddExecutor{barcode: b, validators: validators}
}

func (e *BarcodeAddExecutor) Process(ctx context.Context, req dto.BarcodeAddRequest) (dto.BarcodeResponse, error) {
	if err := processor.Validate(ctx, req, e.validators...); err != nil {
		return dto.BarcodeResponse{}, err
	}
	b, err := e.barcode.Add(ctx, req.ProdukID, domain.BarcodeInput{
		Kode: req.Kode, IsiPerKemasan: req.IsiPerKemasan, Sumber: req.Sumber,
	})
	if err != nil {
		return dto.BarcodeResponse{}, err
	}
	return mapper.ToBarcodeResponse(b), nil
}
