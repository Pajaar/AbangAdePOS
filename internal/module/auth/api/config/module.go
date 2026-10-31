// Package config merakit modul auth: repository → inservice → executor →
// controller, lalu mendaftarkan route-nya (setara ModuleConfig,
// ExecutorConfiguration, dan ControllerConfig di primecash; di Go dirakit
// manual karena tidak ada component scan).
package config

import (
	"database/sql"

	"pos-jamu-backend/internal/common/route"
	appconfig "pos-jamu-backend/internal/config"
	"pos-jamu-backend/internal/middleware"
	"pos-jamu-backend/internal/module/auth/api/controller"
	"pos-jamu-backend/internal/module/auth/api/executor"
	"pos-jamu-backend/internal/module/auth/service/inservice"
	tck "pos-jamu-backend/internal/module/auth/tck/inservice"
)

type Module struct {
	cfg        *appconfig.Config
	controller *controller.AuthController

	// PermissionInService dipakai authorize.Store saat startup.
	PermissionInService tck.IPermissionInService
}

var _ route.Registrar = (*Module)(nil)

func NewModule(cfg *appconfig.Config, db *sql.DB) (*Module, error) {
	userInService, err := inservice.NewUserInService(db)
	if err != nil {
		return nil, err
	}
	sessionInService := inservice.NewSessionInService(db, cfg.RefreshTokenTTL)
	menuInService := inservice.NewMenuInService(db)

	at := executor.AccessTokenConfig{Secret: cfg.JWTSecret, TTL: cfg.AccessTokenTTL}

	ctrl := controller.NewAuthController(controller.Executors{
		Login:   executor.NewLoginExecutor(userInService, sessionInService, at),
		Refresh: executor.NewRefreshExecutor(sessionInService, at),
		Logout:  executor.NewLogoutExecutor(sessionInService),
		Me:      executor.NewMeExecutor(userInService),
		Menus:   executor.NewMenusExecutor(menuInService),
	}, controller.CookieConfig{
		Name:     cfg.CookieName,
		Secure:   cfg.CookieSecure,
		SameSite: cfg.CookieSameSite,
		Domain:   cfg.CookieDomain,
	})

	return &Module{
		cfg:                 cfg,
		controller:          ctrl,
		PermissionInService: inservice.NewPermissionInService(db),
	}, nil
}

// RegisterRoutes: endpoint /auth/* tidak lewat Authorize (lihat CLAUDE.md).
func (m *Module) RegisterRoutes(g route.Groups) {
	loginLimit := middleware.RateLimit(m.cfg.RateLimitLoginPerMinute, m.cfg.RateLimitLoginBurst)

	g.Public.POST("/auth/login", loginLimit, m.controller.Login)
	g.Public.POST("/auth/refresh", m.controller.Refresh)
	g.Public.POST("/auth/logout", m.controller.Logout)

	g.Authenticated.GET("/auth/me", m.controller.Me)
	g.Authenticated.GET("/auth/menus", m.controller.Menus)
}
