// Package repository membungkus query sqlc modul produk dan memetakan
// hasilnya ke tipe domain. Konstruktor menerima sqlcgen.DBTX supaya bisa
// dipakai dengan *sql.DB maupun *sql.Tx.
package repository

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

func nullString(p *string) sql.NullString {
	if p == nil || *p == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: *p, Valid: true}
}

func strPtr(n sql.NullString) *string {
	if !n.Valid {
		return nil
	}
	s := n.String
	return &s
}

func nullDecimal(p *decimal.Decimal) decimal.NullDecimal {
	if p == nil {
		return decimal.NullDecimal{}
	}
	return decimal.NullDecimal{Decimal: *p, Valid: true}
}

func decPtr(n decimal.NullDecimal) *decimal.Decimal {
	if !n.Valid {
		return nil
	}
	d := n.Decimal
	return &d
}
