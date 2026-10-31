package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func server(t *testing.T, status int, body string, delay time.Duration) *httptest.Server {
	t.Helper()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if delay > 0 {
			select {
			case <-time.After(delay):
			case <-r.Context().Done():
				return
			}
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(s.Close)
	return s
}

func TestOpenFoodFacts(t *testing.T) {
	ctx := context.Background()

	found := server(t, 200, `{"status":1,"product":{"product_name":" Tolak Angin ","brands":"Sido Muncul, Lain"}}`, 0)
	p, err := NewOpenFoodFacts(found.URL).Lookup(ctx, "8992761136109")
	if err != nil || !p.Found || p.Nama != "Tolak Angin" || p.Merek != "Sido Muncul" || len(p.Raw) == 0 {
		t.Errorf("found: %+v %v", p, err)
	}

	notFound := server(t, 404, `{"status":0,"status_verbose":"product not found"}`, 0)
	p, err = NewOpenFoodFacts(notFound.URL).Lookup(ctx, "1")
	if err != nil || p.Found {
		t.Errorf("not found harus Found=false tanpa error: %+v %v", p, err)
	}

	broken := server(t, 503, `oops`, 0)
	if _, err := NewOpenFoodFacts(broken.URL).Lookup(ctx, "1"); err == nil {
		t.Error("503 harus error")
	}
}

func TestUPCItemDB(t *testing.T) {
	ctx := context.Background()

	found := server(t, 200, `{"code":"OK","total":1,"items":[{"title":"Minyak Kayu Putih","brand":"Cap Lang"}]}`, 0)
	p, err := NewUPCItemDB(found.URL).Lookup(ctx, "1")
	if err != nil || !p.Found || p.Nama != "Minyak Kayu Putih" || p.Merek != "Cap Lang" {
		t.Errorf("found: %+v %v", p, err)
	}

	empty := server(t, 200, `{"code":"OK","total":0,"items":[]}`, 0)
	if p, err := NewUPCItemDB(empty.URL).Lookup(ctx, "1"); err != nil || p.Found {
		t.Errorf("total 0 = tidak ditemukan: %+v %v", p, err)
	}

	limited := server(t, 429, `{"code":"TOO_FAST"}`, 0)
	if _, err := NewUPCItemDB(limited.URL).Lookup(ctx, "1"); err == nil {
		t.Error("429 harus error (tidak di-cache)")
	}
}

func TestLookupTimeout(t *testing.T) {
	slow := server(t, 200, `{"status":1}`, LookupTimeout+2*time.Second)
	start := time.Now()
	_, err := NewOpenFoodFacts(slow.URL).Lookup(context.Background(), "1")
	if err == nil || !strings.Contains(err.Error(), "deadline") {
		t.Errorf("harus timeout, got %v", err)
	}
	if time.Since(start) > LookupTimeout+time.Second {
		t.Errorf("timeout terlalu lama: %s", time.Since(start))
	}
}
