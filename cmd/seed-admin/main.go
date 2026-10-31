// cmd/seed-admin membuat user pertama (biasanya headadmin) dengan password
// yang sudah di-hash bcrypt. Dijalankan manual sekali saat setup awal,
// bukan lewat endpoint API.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"pos-jamu-backend/internal/config"
	"pos-jamu-backend/internal/database"
	"pos-jamu-backend/internal/module/auth/service/inservice"
	"pos-jamu-backend/internal/module/auth/tck/domain"
)

func main() {
	username := flag.String("username", "", "username login (wajib)")
	password := flag.String("password", "", "password awal (wajib)")
	nama := flag.String("nama", "", "nama lengkap user (wajib)")
	roleKode := flag.String("role", "headadmin", "kode role, lihat tabel roles")
	flag.Parse()

	if *username == "" || *password == "" || *nama == "" {
		fmt.Fprintln(os.Stderr, "usage: seed-admin -username=<u> -password=<p> -nama=<n> [-role=headadmin]")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.NewMySQL(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := inservice.NewUserInService(db)
	if err != nil {
		log.Fatalf("user service: %v", err)
	}

	// Cek role ada, username belum dipakai, dan hash bcrypt dilakukan di service.
	id, err := users.Create(ctx, domain.NewUser{
		Nama:     *nama,
		Username: *username,
		Password: *password,
		RoleKode: *roleKode,
	})
	if err != nil {
		log.Fatalf("buat user: %v", err)
	}

	log.Printf("user %q (id=%d, role=%s) berhasil dibuat", *username, id, *roleKode)
}
