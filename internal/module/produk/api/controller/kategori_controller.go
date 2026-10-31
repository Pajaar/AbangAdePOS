package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pos-jamu-backend/internal/common/httpctx"
	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/util"
)

type KategoriExecutors struct {
	List   processor.IExecutorProcessor[dto.KategoriListRequest, []dto.KategoriResponse]
	Create processor.IExecutorProcessor[dto.KategoriRequest, dto.KategoriResponse]
	Update processor.IExecutorProcessor[dto.KategoriRequest, dto.KategoriResponse]
}

type KategoriController struct {
	ex KategoriExecutors
}

func NewKategoriController(ex KategoriExecutors) *KategoriController {
	return &KategoriController{ex: ex}
}

// GET /api/v1/kategori?aktif=true
func (h *KategoriController) List(c *gin.Context) {
	var req dto.KategoriListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		util.Error(c, httpctx.BindError(err))
		return
	}
	res, err := h.ex.List.Process(c.Request.Context(), req)
	respond(c, http.StatusOK, res, err)
}

// POST /api/v1/kategori
func (h *KategoriController) Create(c *gin.Context) {
	var req dto.KategoriRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, httpctx.BindError(err))
		return
	}
	res, err := h.ex.Create.Process(c.Request.Context(), req)
	respond(c, http.StatusCreated, res, err)
}

// PUT /api/v1/kategori/:id
func (h *KategoriController) Update(c *gin.Context) {
	id, err := httpctx.ParamID(c, "id")
	if err != nil {
		util.Error(c, err)
		return
	}
	var req dto.KategoriRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, httpctx.BindError(err))
		return
	}
	req.ID = id
	res, err := h.ex.Update.Process(c.Request.Context(), req)
	respond(c, http.StatusOK, res, err)
}
