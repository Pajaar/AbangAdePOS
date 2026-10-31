package httpctx

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/util/apperror"
)

var registerOnce sync.Once

// UseJSONFieldNames membuat error validasi binding memakai nama field JSON
// (harga_normal), bukan nama field Go (HargaNormal). Dipanggil sekali saat
// router dibuat.
func UseJSONFieldNames() {
	registerOnce.Do(func() {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			v.RegisterTagNameFunc(func(f reflect.StructField) string {
				for _, tag := range []string{"json", "form"} {
					name := strings.Split(f.Tag.Get(tag), ",")[0]
					if name != "" && name != "-" {
						return name
					}
				}
				return f.Name
			})
		}
	})
}

// BindError mengubah error ShouldBind* menjadi apperror 400. Error validasi
// tag (required, max, oneof, ...) dikembalikan sebagai daftar details.
func BindError(err error) error {
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		details := make([]processor.ValidationError, 0, len(verrs))
		for _, fe := range verrs {
			details = append(details, processor.ValidationError{
				Field:   fieldPath(fe),
				Code:    fe.Tag(),
				Message: pesanTag(fe),
			})
		}
		return apperror.Validation(details)
	}
	return apperror.BadRequest("Format data tidak valid")
}

// fieldPath: "ProdukRequest.barcodes[0].kode" -> "barcodes[0].kode".
// Segmen nama tipe Go (huruf besar, misalnya struct yang di-embed) dibuang.
func fieldPath(fe validator.FieldError) string {
	parts := strings.Split(fe.Namespace(), ".")
	out := make([]string, 0, len(parts))
	for _, p := range parts[1:] {
		if p != "" && p[0] >= 'A' && p[0] <= 'Z' {
			continue
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		return fe.Field()
	}
	return strings.Join(out, ".")
}

func pesanTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Wajib diisi"
	case "max":
		return "Maksimal " + fe.Param()
	case "min":
		return "Minimal " + fe.Param()
	case "oneof":
		return "Harus salah satu dari: " + fe.Param()
	}
	return "Tidak valid"
}

// ParamID membaca parameter path angka positif, misalnya :id.
func ParamID(c *gin.Context, name string) (uint64, error) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || id == 0 {
		return 0, apperror.BadRequest("Parameter " + name + " tidak valid")
	}
	return id, nil
}
