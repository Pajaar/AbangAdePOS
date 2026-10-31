package processor

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"pos-jamu-backend/internal/util/apperror"
)

type req struct{ Kode string }

type kodeWajibValidator struct{}

func (kodeWajibValidator) Process(_ context.Context, r req) ([]ValidationError, error) {
	if r.Kode == "" {
		return OneError("kode", "required", "Kode wajib diisi"), nil
	}
	return NoErrors(), nil
}

type gagalValidator struct{}

func (gagalValidator) Process(context.Context, req) ([]ValidationError, error) {
	return nil, errors.New("db down")
}

func TestValidate(t *testing.T) {
	ctx := context.Background()

	if err := Validate[req](ctx, req{Kode: "A"}, kodeWajibValidator{}); err != nil {
		t.Errorf("valid request seharusnya lolos: %v", err)
	}

	err := Validate[req](ctx, req{}, kodeWajibValidator{}, kodeWajibValidator{})
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) || appErr.Status != http.StatusBadRequest {
		t.Fatalf("expected 400, got %v", err)
	}
	if details, ok := appErr.Details.([]ValidationError); !ok || len(details) != 2 {
		t.Errorf("details harus berisi 2 error, got %#v", appErr.Details)
	}

	if err := Validate[req](ctx, req{}, gagalValidator{}); err == nil || errors.As(err, &appErr) {
		t.Errorf("error validator harus diteruskan apa adanya, got %v", err)
	}
}
