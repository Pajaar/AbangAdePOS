package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pos-jamu-backend/internal/common/httpctx"
	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/util"
)

type BarcodeExecutors struct {
	Add    processor.IExecutorProcessor[dto.BarcodeAddRequest, dto.BarcodeResponse]
	Delete processor.IExecutorProcessor[dto.IDRequest, processor.Empty]
	Kasir  processor.IExecutorProcessor[dto.BarcodeKasirRequest, dto.BarcodeKasirResponse]
	Lookup processor.IExecutorProcessor[dto.BarcodeLookupRequest, dto.BarcodeLookupResponse]
}

type BarcodeController struct {
	ex BarcodeExecutors
}

func NewBarcodeController(ex BarcodeExecutors) *BarcodeController {
	return &BarcodeController{ex: ex}
}

// POST /api/v1/produk/:id/barcode
func (h *BarcodeController) Add(c *gin.Context) {
	id, err := httpctx.ParamID(c, "id")
	if err != nil {
		util.Error(c, err)
		return
	}
	var req dto.BarcodeAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, httpctx.BindError(err))
		return
	}
	req.ProdukID = id
	res, err := h.ex.Add.Process(c.Request.Context(), req)
	respond(c, http.StatusCreated, res, err)
}

// DELETE /api/v1/produk-barcode/:id
func (h *BarcodeController) Delete(c *gin.Context) {
	id, err := httpctx.ParamID(c, "id")
	if err != nil {
		util.Error(c, err)
		return
	}
	res, err := h.ex.Delete.Process(c.Request.Context(), dto.IDRequest{ID: id})
	respond(c, http.StatusOK, res, err)
}

// GET /api/v1/barcode/:kode (kasir: 404 tidak terdaftar, 409 produk draft)
func (h *BarcodeController) Kasir(c *gin.Context) {
	role, _ := httpctx.RoleKode(c)
	res, err := h.ex.Kasir.Process(c.Request.Context(), dto.BarcodeKasirRequest{Kode: c.Param("kode"), RoleKode: role})
	respond(c, http.StatusOK, res, err)
}

// POST /api/v1/barcode/lookup
func (h *BarcodeController) Lookup(c *gin.Context) {
	var req dto.BarcodeLookupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, httpctx.BindError(err))
		return
	}
	req.RoleKode, _ = httpctx.RoleKode(c)
	res, err := h.ex.Lookup.Process(c.Request.Context(), req)
	respond(c, http.StatusOK, res, err)
}
