package util

import (
	"errors"
	"strings"
)

// Jenis barcode, sama dengan ENUM produk_barcode.jenis.
const (
	BarcodeEAN13    = "EAN13"
	BarcodeEAN8     = "EAN8"
	BarcodeUPCA     = "UPCA"
	BarcodeInternal = "INTERNAL"
)

var (
	ErrBarcodeKosong   = errors.New("barcode kosong")
	ErrBarcodeChecksum = errors.New("digit kontrol barcode tidak valid")
	ErrBarcodeTidakSah = errors.New("barcode hanya boleh huruf dan angka, maksimal 14 karakter")
)

// NormalizeBarcode menyiapkan barcode sebelum disimpan atau dicari.
//   - 12 digit (UPC-A) diubah ke EAN-13 dengan awalan 0, jenis UPCA.
//   - 8 dan 13 digit divalidasi digit kontrolnya (EAN-8 / EAN-13).
//   - Selain itu dianggap kode INTERNAL toko (huruf/angka, maks 14 karakter).
//
// Spasi dan tanda hubung di awal/tengah diabaikan.
func NormalizeBarcode(raw string) (kode, jenis string, err error) {
	s := strings.ToUpper(strings.NewReplacer(" ", "", "-", "").Replace(strings.TrimSpace(raw)))
	if s == "" {
		return "", "", ErrBarcodeKosong
	}

	if isDigits(s) {
		switch len(s) {
		case 8:
			if !validGTIN(s) {
				return "", "", ErrBarcodeChecksum
			}
			return s, BarcodeEAN8, nil
		case 12:
			ean := "0" + s
			if !validGTIN(ean) {
				return "", "", ErrBarcodeChecksum
			}
			return ean, BarcodeUPCA, nil
		case 13:
			if !validGTIN(s) {
				return "", "", ErrBarcodeChecksum
			}
			return s, BarcodeEAN13, nil
		}
	}

	if len(s) > 14 || !isAlnum(s) {
		return "", "", ErrBarcodeTidakSah
	}
	return s, BarcodeInternal, nil
}

// validGTIN memeriksa digit kontrol EAN-8/EAN-13 (bobot 3 dan 1 dari kanan).
func validGTIN(s string) bool {
	sum := 0
	n := len(s)
	for i := 0; i < n-1; i++ {
		d := int(s[i] - '0')
		if (n-1-i)%2 == 1 {
			sum += d * 3
		} else {
			sum += d
		}
	}
	check := (10 - sum%10) % 10
	return check == int(s[n-1]-'0')
}

func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func isAlnum(s string) bool {
	for _, r := range s {
		if (r < '0' || r > '9') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}
