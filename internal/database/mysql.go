package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"pos-jamu-backend/internal/config"
)

// NewMySQL membuka connection pool ke MySQL (Railway atau lokal).
// ConnMaxLifetime dibatasi karena connection pooler di Railway bisa
// menutup koneksi idle secara sepihak.
func NewMySQL(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.MySQLDSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
