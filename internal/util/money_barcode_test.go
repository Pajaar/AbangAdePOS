package util

import (
	"testing"

	"github.com/shopspring/decimal"
)

func d(s string) decimal.Decimal { return decimal.RequireFromString(s) }

func TestHPPRataRataBaru(t *testing.T) {
	cases := []struct {
		name     string
		stok     int64
		hpp      string
		qty      int64
		harga    string
		expected string
	}{
		{"stok kosong pakai harga beli", 0, "0", 10, "2000", "2000"},
		{"stok minus pakai harga beli", -3, "1500", 10, "2000", "2000"},
		{"rata-rata", 10, "2000", 10, "3000", "2500"},
		{"dibulatkan 2 desimal", 3, "1000", 1, "2000", "1250"},
		{"pecahan", 2, "1000", 1, "1001", "1000.33"},
	}
	for _, c := range cases {
		got := HPPRataRataBaru(c.stok, d(c.hpp), c.qty, d(c.harga))
		if !got.Equal(d(c.expected)) {
			t.Errorf("%s: got %s, want %s", c.name, got, c.expected)
		}
	}
}

func TestHargaNettTengah(t *testing.T) {
	nett := HargaNett(d("10000"), d("20"))
	if !nett.Equal(d("12000")) {
		t.Errorf("nett = %s", nett)
	}
	if n := HargaNett(d("10100"), d("20")); !n.Equal(d("12500")) { // 12120 -> 12500
		t.Errorf("nett pembulatan = %s", n)
	}
	if tg := HargaTengah(d("15000"), d("12000")); !tg.Equal(d("13500")) {
		t.Errorf("tengah = %s", tg)
	}
	if tg := HargaTengah(d("15000"), d("12500")); !tg.Equal(d("14000")) { // 13750 -> 14000
		t.Errorf("tengah pembulatan = %s", tg)
	}
}

func TestNormalizeBarcode(t *testing.T) {
	cases := []struct {
		in, kode, jenis string
		err             error
	}{
		{"4006381333931", "4006381333931", BarcodeEAN13, nil},
		{"036000291452", "0036000291452", BarcodeUPCA, nil},
		{"0 3600-0291452", "0036000291452", BarcodeUPCA, nil},
		{"96385074", "96385074", BarcodeEAN8, nil},
		{"4006381333932", "", "", ErrBarcodeChecksum},
		{"jmu-01", "JMU01", BarcodeInternal, nil},
		{"   ", "", "", ErrBarcodeKosong},
		{"ABC_123", "", "", ErrBarcodeTidakSah},
		{"123456789012345", "", "", ErrBarcodeTidakSah},
	}
	for _, c := range cases {
		kode, jenis, err := NormalizeBarcode(c.in)
		if err != c.err || kode != c.kode || jenis != c.jenis {
			t.Errorf("%q: got (%q,%q,%v), want (%q,%q,%v)", c.in, kode, jenis, err, c.kode, c.jenis, c.err)
		}
	}
}
