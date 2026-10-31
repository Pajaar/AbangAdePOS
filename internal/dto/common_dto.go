package dto

// PagedResponse dipakai semua endpoint daftar yang memakai paging.
type PagedResponse[T any] struct {
	Items  []T   `json:"items"`
	Total  int64 `json:"total"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

// IDRequest untuk endpoint yang hanya butuh :id dari path.
type IDRequest struct {
	ID       uint64
	RoleKode string
}
