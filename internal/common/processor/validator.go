package processor

import (
	"context"

	"pos-jamu-backend/internal/util/apperror"
)

// IValidatorProcessor memeriksa satu aturan bisnis (biasanya yang butuh
// database, misalnya "kode sudah dipakai"). Validasi format field cukup
// lewat tag binding di DTO.
type IValidatorProcessor[Req any] interface {
	Process(ctx context.Context, req Req) ([]ValidationError, error)
}

type ValidationError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NoErrors() []ValidationError {
	return nil
}

func OneError(field, code, message string) []ValidationError {
	return []ValidationError{{Field: field, Code: code, Message: message}}
}

// Validate menjalankan semua validator dan mengumpulkan errornya, setara
// ValidationDetailsReport di primecash. Kembali apperror.Validation (400)
// bila ada pelanggaran, atau error asli bila validator gagal dijalankan.
func Validate[Req any](ctx context.Context, req Req, validators ...IValidatorProcessor[Req]) error {
	var all []ValidationError
	for _, v := range validators {
		errs, err := v.Process(ctx, req)
		if err != nil {
			return err
		}
		all = append(all, errs...)
	}
	if len(all) > 0 {
		return apperror.Validation(all)
	}
	return nil
}
