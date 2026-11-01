// Package controller berisi handler HTTP modul menu jamu.
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pos-jamu-backend/internal/common/httpctx"
	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/util"
)

type Executors struct {
	JamuList         processor.IExecutorProcessor[dto.JamuListRequest, []dto.JamuResponse]
	JamuCreate       processor.IExecutorProcessor[dto.JamuRequest, dto.JamuResponse]
	JamuUpdate       processor.IExecutorProcessor[dto.JamuRequest, dto.JamuResponse]
	AdditionalList   processor.IExecutorProcessor[dto.AdditionalListRequest, []dto.AdditionalResponse]
	AdditionalCreate processor.IExecutorProcessor[dto.AdditionalRequest, dto.AdditionalResponse]
	AdditionalUpdate processor.IExecutorProcessor[dto.AdditionalRequest, dto.AdditionalResponse]
}

type JamuController struct {
	ex Executors
}

func NewJamuController(ex Executors) *JamuController {
	return &JamuController{ex: ex}
}

func respond[T any](c *gin.Context, status int, res T, err error) {
	if err != nil {
		util.Error(c, err)
		return
	}
	util.Success(c, status, res)
}

// GET /api/v1/jamu?aktif=
func (h *JamuController) JamuList(c *gin.Context) {
	var req dto.JamuListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		util.Error(c, httpctx.BindError(err))
		return
	}
	req.RoleKode, _ = httpctx.RoleKode(c)
	res, err := h.ex.JamuList.Process(c.Request.Context(), req)
	respond(c, http.StatusOK, res, err)
}

// POST /api/v1/jamu
func (h *JamuController) JamuCreate(c *gin.Context) {
	var req dto.JamuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, httpctx.BindError(err))
		return
	}
	req.RoleKode, _ = httpctx.RoleKode(c)
	res, err := h.ex.JamuCreate.Process(c.Request.Context(), req)
	respond(c, http.StatusCreated, res, err)
}

// PUT /api/v1/jamu/:id
func (h *JamuController) JamuUpdate(c *gin.Context) {
	id, err := httpctx.ParamID(c, "id")
	if err != nil {
		util.Error(c, err)
		return
	}
	var req dto.JamuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, httpctx.BindError(err))
		return
	}
	req.ID = id
	req.RoleKode, _ = httpctx.RoleKode(c)
	res, err := h.ex.JamuUpdate.Process(c.Request.Context(), req)
	respond(c, http.StatusOK, res, err)
}

// GET /api/v1/additional-items?aktif=
func (h *JamuController) AdditionalList(c *gin.Context) {
	var req dto.AdditionalListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		util.Error(c, httpctx.BindError(err))
		return
	}
	req.RoleKode, _ = httpctx.RoleKode(c)
	res, err := h.ex.AdditionalList.Process(c.Request.Context(), req)
	respond(c, http.StatusOK, res, err)
}

// POST /api/v1/additional-items
func (h *JamuController) AdditionalCreate(c *gin.Context) {
	var req dto.AdditionalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, httpctx.BindError(err))
		return
	}
	req.RoleKode, _ = httpctx.RoleKode(c)
	res, err := h.ex.AdditionalCreate.Process(c.Request.Context(), req)
	respond(c, http.StatusCreated, res, err)
}

// PUT /api/v1/additional-items/:id
func (h *JamuController) AdditionalUpdate(c *gin.Context) {
	id, err := httpctx.ParamID(c, "id")
	if err != nil {
		util.Error(c, err)
		return
	}
	var req dto.AdditionalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, httpctx.BindError(err))
		return
	}
	req.ID = id
	req.RoleKode, _ = httpctx.RoleKode(c)
	res, err := h.ex.AdditionalUpdate.Process(c.Request.Context(), req)
	respond(c, http.StatusOK, res, err)
}
