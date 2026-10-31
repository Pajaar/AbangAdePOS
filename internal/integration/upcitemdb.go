package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"pos-jamu-backend/internal/module/produk/tck/domain"
	"pos-jamu-backend/internal/module/produk/tck/inservice"
)

// UPCItemDB memakai endpoint trial (gratis, dibatasi sekitar 100 request/hari).
type UPCItemDB struct {
	baseURL string
	client  *http.Client
}

var _ inservice.IBarcodeSource = (*UPCItemDB)(nil)

// NewUPCItemDB: baseURL kosong = https://api.upcitemdb.com
func NewUPCItemDB(baseURL string) *UPCItemDB {
	if baseURL == "" {
		baseURL = "https://api.upcitemdb.com"
	}
	return &UPCItemDB{baseURL: strings.TrimRight(baseURL, "/"), client: &http.Client{}}
}

func (u *UPCItemDB) Name() string { return "upcitemdb" }

func (u *UPCItemDB) Lookup(ctx context.Context, kode string) (domain.ExternalProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, LookupTimeout)
	defer cancel()

	body, status, err := get(ctx, u.client, u.baseURL+"/prod/trial/lookup?upc="+url.QueryEscape(kode))
	if err != nil {
		return domain.ExternalProduct{}, err
	}

	// 404 / 400 INVALID_UPC = barcode tidak ada di database mereka.
	if status == http.StatusNotFound || status == http.StatusBadRequest {
		return domain.ExternalProduct{Found: false}, nil
	}
	// 429 (batas harian/terlalu cepat) dan 5xx = error, jangan di-cache.
	if status != http.StatusOK {
		return domain.ExternalProduct{}, fmt.Errorf("upcitemdb: HTTP %d", status)
	}

	var resp struct {
		Code  string `json:"code"`
		Total int    `json:"total"`
		Items []struct {
			Title string `json:"title"`
			Brand string `json:"brand"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return domain.ExternalProduct{}, fmt.Errorf("upcitemdb: respons tidak valid: %w", err)
	}
	if resp.Code != "OK" {
		return domain.ExternalProduct{}, fmt.Errorf("upcitemdb: code %s", resp.Code)
	}
	if resp.Total == 0 || len(resp.Items) == 0 {
		return domain.ExternalProduct{Found: false}, nil
	}
	return domain.ExternalProduct{
		Found: true,
		Nama:  strings.TrimSpace(resp.Items[0].Title),
		Merek: strings.TrimSpace(resp.Items[0].Brand),
		Raw:   body,
	}, nil
}
