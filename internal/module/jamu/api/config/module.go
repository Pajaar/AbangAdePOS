// Package config merakit modul menu jamu dan mendaftarkan route-nya.
package config

import (
	"database/sql"

	"pos-jamu-backend/internal/common/access"
	"pos-jamu-backend/internal/common/route"
	"pos-jamu-backend/internal/module/jamu/api/controller"
	"pos-jamu-backend/internal/module/jamu/api/executor"
	"pos-jamu-backend/internal/module/jamu/api/validator"
	"pos-jamu-backend/internal/module/jamu/service/inservice"
	produktck "pos-jamu-backend/internal/module/produk/tck/inservice"
)

type Module struct {
	controller *controller.JamuController
}

var _ route.Registrar = (*Module)(nil)

// NewModule: produk dipakai validator untuk mengecek sachet/produk yang dipilih.
func NewModule(db *sql.DB, actions access.IActionChecker, produk produktck.IProdukInService) *Module {
	jamuSvc := inservice.NewJamuInService(db)
	additionalSvc := inservice.NewAdditionalItemInService(db)

	jamuValidator := validator.NewJamuValidator(produk)
	additionalValidator := validator.NewAdditionalValidator(produk)

	return &Module{controller: controller.NewJamuController(controller.Executors{
		JamuList:         executor.NewJamuListExecutor(jamuSvc, actions),
		JamuCreate:       executor.NewJamuCreateExecutor(jamuSvc, actions, jamuValidator),
		JamuUpdate:       executor.NewJamuUpdateExecutor(jamuSvc, actions, jamuValidator),
		AdditionalList:   executor.NewAdditionalListExecutor(additionalSvc, actions),
		AdditionalCreate: executor.NewAdditionalCreateExecutor(additionalSvc, actions, additionalValidator),
		AdditionalUpdate: executor.NewAdditionalUpdateExecutor(additionalSvc, actions, additionalValidator),
	})}
}

// RegisterRoutes: pola path sama dengan api_endpoints di seed 000002 (menu_jamu).
func (m *Module) RegisterRoutes(g route.Groups) {
	p := g.Protected
	p.GET("/jamu", m.controller.JamuList)
	p.POST("/jamu", m.controller.JamuCreate)
	p.PUT("/jamu/:id", m.controller.JamuUpdate)
	p.GET("/additional-items", m.controller.AdditionalList)
	p.POST("/additional-items", m.controller.AdditionalCreate)
	p.PUT("/additional-items/:id", m.controller.AdditionalUpdate)
}
