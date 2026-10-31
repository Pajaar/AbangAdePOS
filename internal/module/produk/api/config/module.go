// Package config merakit modul produk: repository -> inservice -> validator
// -> executor -> controller, lalu mendaftarkan route-nya.
package config

import (
	"database/sql"

	"pos-jamu-backend/internal/common/access"
	"pos-jamu-backend/internal/common/route"
	"pos-jamu-backend/internal/integration"
	jurnaltck "pos-jamu-backend/internal/module/jurnal/tck/inservice"
	"pos-jamu-backend/internal/module/produk/api/controller"
	"pos-jamu-backend/internal/module/produk/api/executor"
	"pos-jamu-backend/internal/module/produk/api/validator"
	"pos-jamu-backend/internal/module/produk/service/inservice"
)

type Module struct {
	kategori *controller.KategoriController
	produk   *controller.ProdukController
	barcode  *controller.BarcodeController
}

var _ route.Registrar = (*Module)(nil)

// NewModule: actions dipakai untuk izin lihat_hpp, jurnal untuk posting
// jurnal penyesuaian stok.
func NewModule(db *sql.DB, actions access.IActionChecker, jurnal jurnaltck.IJurnalInService) *Module {
	kategoriSvc := inservice.NewKategoriInService(db)
	produkSvc := inservice.NewProdukInService(db)
	barcodeSvc := inservice.NewBarcodeInService(db, produkSvc,
		integration.NewOpenFoodFacts(""),
		integration.NewUPCItemDB(""),
	)
	penyesuaianSvc := inservice.NewPenyesuaianStokInService(db, jurnal)

	kategoriNamaUnik := validator.NewKategoriNamaUnikValidator(kategoriSvc)
	produkKategori := validator.NewProdukKategoriValidator(kategoriSvc)
	produkHarga := validator.ProdukHargaValidator{}

	return &Module{
		kategori: controller.NewKategoriController(controller.KategoriExecutors{
			List:   executor.NewKategoriListExecutor(kategoriSvc),
			Create: executor.NewKategoriCreateExecutor(kategoriSvc, kategoriNamaUnik),
			Update: executor.NewKategoriUpdateExecutor(kategoriSvc, kategoriNamaUnik),
		}),
		produk: controller.NewProdukController(controller.ProdukExecutors{
			List:   executor.NewProdukListExecutor(produkSvc, actions),
			Detail: executor.NewProdukDetailExecutor(produkSvc, actions),
			Create: executor.NewProdukCreateExecutor(produkSvc, actions,
				produkKategori, produkHarga, validator.NewProdukBarcodeBelumDipakaiValidator(barcodeSvc)),
			Update:            executor.NewProdukUpdateExecutor(produkSvc, actions, produkKategori, produkHarga),
			PenyesuaianCreate: executor.NewPenyesuaianCreateExecutor(penyesuaianSvc, actions),
			PenyesuaianList:   executor.NewPenyesuaianListExecutor(penyesuaianSvc, produkSvc, actions),
		}),
		barcode: controller.NewBarcodeController(controller.BarcodeExecutors{
			Add:    executor.NewBarcodeAddExecutor(barcodeSvc, validator.NewBarcodeBelumDipakaiValidator(barcodeSvc)),
			Delete: executor.NewBarcodeDeleteExecutor(barcodeSvc),
			Kasir:  executor.NewBarcodeKasirExecutor(barcodeSvc, actions),
			Lookup: executor.NewBarcodeLookupExecutor(barcodeSvc, actions),
		}),
	}
}

// RegisterRoutes: semua route lewat Auth + Authorize. Pola path harus sama
// persis dengan api_endpoints di seed (000002 dan 000005).
func (m *Module) RegisterRoutes(g route.Groups) {
	p := g.Protected

	p.GET("/kategori", m.kategori.List)
	p.POST("/kategori", m.kategori.Create)
	p.PUT("/kategori/:id", m.kategori.Update)

	p.GET("/produk", m.produk.List)
	p.POST("/produk", m.produk.Create)
	p.GET("/produk/:id", m.produk.Detail)
	p.PUT("/produk/:id", m.produk.Update)
	p.POST("/produk/:id/penyesuaian-stok", m.produk.PenyesuaianCreate)
	p.GET("/produk/:id/penyesuaian-stok", m.produk.PenyesuaianList)

	p.POST("/produk/:id/barcode", m.barcode.Add)
	p.DELETE("/produk-barcode/:id", m.barcode.Delete)
	p.GET("/barcode/:kode", m.barcode.Kasir)
	p.POST("/barcode/lookup", m.barcode.Lookup)
}
