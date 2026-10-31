package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pos-jamu-backend/internal/common/httpctx"
	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/util"
	"pos-jamu-backend/internal/util/apperror"
)

type ProdukExecutors struct {
	List              processor.IExecutorProcessor[dto.ProdukListRequest, dto.PagedResponse[dto.ProdukResponse]]
	Detail            processor.IExecutorProcessor[dto.IDRequest, dto.ProdukResponse]
	Create            processor.IExecutorProcessor[dto.ProdukRequest, dto.ProdukResponse]
	Update            processor.IExecutorProcessor[dto.ProdukRequest, dto.ProdukResponse]
	PenyesuaianCreate processor.IExecutorProcessor[dto.PenyesuaianRequest, dto.PenyesuaianResponse]
	PenyesuaianList   processor.IExecutorProcessor[dto.PenyesuaianListRequest, []dto.PenyesuaianResponse]
}

type ProdukController struct {
	ex ProdukExecutors
}

func NewProdukController(ex ProdukExecutors) *ProdukController {
	return &ProdukController{ex: ex}
}

// GET /api/v1/produk?q=&kategori_id=&status=&aktif=&limit=&offset=
func (h *ProdukController) List(c *gin.Context) {
	var req dto.ProdukListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		util.Error(c, httpctx.BindError(err))
		return
	}
	req.RoleKode, _ = httpctx.RoleKode(c)
	res, err := h.ex.List.Process(c.Request.Context(), req)
	respond(c, http.StatusOK, res, err)
}

// GET /api/v1/produk/:id
func (h *ProdukController) Detail(c *gin.Context) {
	id, err := httpctx.ParamID(c, "id")
	if err != nil {
		util.Error(c, err)
		return
	}
	role, _ := httpctx.RoleKode(c)
	res, err := h.ex.Detail.Process(c.Request.Context(), dto.IDRequest{ID: id, RoleKode: role})
	respond(c, http.StatusOK, res, err)
}

// POST /api/v1/produk
func (h *ProdukController) Create(c *gin.Context) {
	var req dto.ProdukRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, httpctx.BindError(err))
		return
	}
	req.RoleKode, _ = httpctx.RoleKode(c)
	res, err := h.ex.Create.Process(c.Request.Context(), req)
	respond(c, http.StatusCreated, res, err)
}

// PUT /api/v1/produk/:id
func (h *ProdukController) Update(c *gin.Context) {
	id, err := httpctx.ParamID(c, "id")
	if err != nil {
		util.Error(c, err)
		return
	}
	var req dto.ProdukRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, httpctx.BindError(err))
		return
	}
	if len(req.Barcodes) > 0 {
		util.Error(c, apperror.BadRequest("Barcode diubah lewat endpoint /produk/:id/barcode, bukan lewat ubah produk"))
		return
	}
	req.ID = id
	req.RoleKode, _ = httpctx.RoleKode(c)
	res, err := h.ex.Update.Process(c.Request.Context(), req)
	respond(c, http.StatusOK, res, err)
}

// POST /api/v1/produk/:id/penyesuaian-stok
func (h *ProdukController) PenyesuaianCreate(c *gin.Context) {
	id, err := httpctx.ParamID(c, "id")
	if err != nil {
		util.Error(c, err)
		return
	}
	userID, ok := httpctx.UserID(c)
	if !ok {
		util.Error(c, apperror.Unauthorized("Belum login"))
		return
	}
	var req dto.PenyesuaianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, httpctx.BindError(err))
		return
	}
	req.ProdukID, req.UserID = id, userID
	req.RoleKode, _ = httpctx.RoleKode(c)
	res, err := h.ex.PenyesuaianCreate.Process(c.Request.Context(), req)
	respond(c, http.StatusCreated, res, err)
}

// GET /api/v1/produk/:id/penyesuaian-stok?limit=
func (h *ProdukController) PenyesuaianList(c *gin.Context) {
	id, err := httpctx.ParamID(c, "id")
	if err != nil {
		util.Error(c, err)
		return
	}
	var req dto.PenyesuaianListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		util.Error(c, httpctx.BindError(err))
		return
	}
	req.ProdukID = id
	req.RoleKode, _ = httpctx.RoleKode(c)
	res, err := h.ex.PenyesuaianList.Process(c.Request.Context(), req)
	respond(c, http.StatusOK, res, err)
}
