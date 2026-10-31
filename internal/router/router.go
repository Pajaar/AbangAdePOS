package router

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"pos-jamu-backend/internal/authorize"
	"pos-jamu-backend/internal/common/route"
	"pos-jamu-backend/internal/config"
	"pos-jamu-backend/internal/middleware"
	"pos-jamu-backend/internal/util"
)

type Deps struct {
	Cfg            *config.Config
	DB             *sql.DB
	AuthorizeStore *authorize.Store
	// Modules: config setiap modul (auth, produk, ...) yang mendaftarkan
	// route-nya sendiri ke grup yang sudah diberi middleware di sini.
	Modules []route.Registrar
}

// New menyusun middleware dan route sesuai urutan wajib di CLAUDE.md:
// Recovery -> RequestID -> Logger -> CORS, baru kemudian route group.
// CORS harus sebelum Auth karena preflight OPTIONS tidak membawa token.
func New(deps Deps) *gin.Engine {
	r := gin.New()

	r.Use(middleware.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	// CORS hanya untuk development (Vite dev server). Di produksi frontend
	// dan API satu origin, CORS_ALLOWED_ORIGINS dikosongkan, dan
	// gin-contrib/cors panic bila dipasang tanpa origin.
	if len(deps.Cfg.CORSAllowedOrigins) > 0 {
		r.Use(middleware.CORS(deps.Cfg.CORSAllowedOrigins))
	}

	r.GET("/healthz", healthzHandler(deps.DB))

	if deps.Cfg.AppEnv == "development" {
		r.GET("/debug/panic", func(c *gin.Context) {
			panic("debug panic test")
		})
	}

	api := r.Group("/api/v1")
	groups := route.Groups{
		Public:        api,
		Authenticated: api.Group("", middleware.Auth(deps.Cfg.JWTSecret)),
		Protected:     api.Group("", middleware.Auth(deps.Cfg.JWTSecret), middleware.Authorize(deps.AuthorizeStore)),
	}
	for _, m := range deps.Modules {
		m.RegisterRoutes(groups)
	}
	protected := groups.Protected

	// Stub sementara untuk memverifikasi Authorize (401/403/200) sebelum
	// fitur Produk sungguhan ada di Tahap 2. Pola endpoint ini sudah
	// terdaftar di seed akses (menu 'produk').
	protected.GET("/kategori", func(c *gin.Context) {
		util.Success(c, http.StatusOK, []any{})
	})

	return r
}

func healthzHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db_unreachable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
