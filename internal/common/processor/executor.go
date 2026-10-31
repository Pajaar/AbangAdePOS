// Package processor berisi kontrak dasar lapisan api, setara
// IExecutorProcessor dan IValidatorProcessor di primecash.
package processor

import "context"

// IExecutorProcessor menjalankan satu use case (satu aksi API).
// Controller hanya mem-bind request lalu memanggil Process.
type IExecutorProcessor[Req any, Res any] interface {
	Process(ctx context.Context, req Req) (Res, error)
}

// Empty dipakai sebagai Res untuk executor yang tidak mengembalikan data.
type Empty struct{}
