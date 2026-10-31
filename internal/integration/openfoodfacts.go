// Package integration berisi klien API pihak ketiga.
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"pos-jamu-backend/internal/module/produk/tck/domain"
	"pos-jamu-backend/internal/module/produk/tck/inservice"
)

// LookupTimeout = batas waktu per API (CLAUDE.md: 3 detik).
const LookupTimeout = 3 * time.Second

const userAgent = "POSJamu/1.0 (warung jamu, lookup barcode)"

type OpenFoodFacts struct {
	baseURL string
	client  *http.Client
}

var _ inservice.IBarcodeSource = (*OpenFoodFacts)(nil)

// NewOpenFoodFacts: baseURL kosong = https://world.openfoodfacts.org
func NewOpenFoodFacts(baseURL string) *OpenFoodFacts {
	if baseURL == "" {
		baseURL = "https://world.openfoodfacts.org"
	}
	return &OpenFoodFacts{baseURL: strings.TrimRight(baseURL, "/"), client: &http.Client{}}
}

func (o *OpenFoodFacts) Name() string { return "openfoodfacts" }

func (o *OpenFoodFacts) Lookup(ctx context.Context, kode string) (domain.ExternalProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, LookupTimeout)
	defer cancel()

	u := fmt.Sprintf("%s/api/v2/product/%s.json?fields=product_name,brands", o.baseURL, url.PathEscape(kode))
	body, status, err := get(ctx, o.client, u)
	if err != nil {
		return domain.ExternalProduct{}, err
	}

	var resp struct {
		Status  int `json:"status"`
		Product struct {
			ProductName string `json:"product_name"`
			Brands      string `json:"brands"`
		} `json:"product"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return domain.ExternalProduct{}, fmt.Errorf("openfoodfacts: respons tidak valid (HTTP %d): %w", status, err)
	}

	switch {
	case status == http.StatusOK && resp.Status == 1:
		return domain.ExternalProduct{
			Found: true,
			Nama:  strings.TrimSpace(resp.Product.ProductName),
			Merek: firstBrand(resp.Product.Brands),
			Raw:   body,
		}, nil
	case (status == http.StatusOK || status == http.StatusNotFound) && resp.Status == 0:
		return domain.ExternalProduct{Found: false}, nil
	}
	return domain.ExternalProduct{}, fmt.Errorf("openfoodfacts: HTTP %d", status)
}

// firstBrand: Open Food Facts mengirim beberapa merek dipisah koma.
func firstBrand(brands string) string {
	return strings.TrimSpace(strings.Split(brands, ",")[0])
}

// get menjalankan GET dan membaca body (maks 1 MB).
func get(ctx context.Context, client *http.Client, u string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return nil, res.StatusCode, err
	}
	return body, res.StatusCode, nil
}
