package main

import (
	"context"
	"log"
	"time"

	"pos-jamu-backend/internal/authorize"
	"pos-jamu-backend/internal/common/route"
	"pos-jamu-backend/internal/config"
	"pos-jamu-backend/internal/database"
	authconfig "pos-jamu-backend/internal/module/auth/api/config"
	"pos-jamu-backend/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.NewMySQL(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := db.PingContext(pingCtx); err != nil {
		pingCancel()
		log.Fatalf("database ping: %v", err)
	}
	pingCancel()

	// Rakit setiap modul. Modul baru (produk, penjualan, ...) ditambahkan di sini.
	authModule, err := authconfig.NewModule(cfg, db)
	if err != nil {
		log.Fatalf("auth module: %v", err)
	}

	authorizeStore := authorize.NewStore()
	loadCtx, loadCancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := authorizeStore.Load(loadCtx, authModule.PermissionInService); err != nil {
		loadCancel()
		log.Fatalf("load permissions: %v", err)
	}
	loadCancel()

	r := router.New(router.Deps{
		Cfg:            cfg,
		DB:             db,
		AuthorizeStore: authorizeStore,
		Modules:        []route.Registrar{authModule},
	})

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server: %v", err)
	}
}
