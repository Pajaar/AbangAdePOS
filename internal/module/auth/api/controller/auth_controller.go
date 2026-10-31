// Package controller berisi handler HTTP modul auth (setara api/controller
// di primecash). Controller hanya: bind/validasi format input, isi data dari
// HTTP (cookie, header, context), panggil executor, lalu kirim respons.
package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"pos-jamu-backend/internal/common/processor"
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/util"
	"pos-jamu-backend/internal/util/apperror"
)

// CookieConfig mengatur cookie refresh token. Path dibatasi /api/v1/auth
// supaya cookie tidak ikut terkirim ke endpoint lain.
type CookieConfig struct {
	Name     string
	Secure   bool
	SameSite string
	Domain   string
}

const refreshCookiePath = "/api/v1/auth"

type AuthController struct {
	loginExecutor   processor.IExecutorProcessor[dto.LoginRequest, dto.LoginResult]
	refreshExecutor processor.IExecutorProcessor[dto.RefreshRequest, dto.RefreshResult]
	logoutExecutor  processor.IExecutorProcessor[dto.LogoutRequest, processor.Empty]
	meExecutor      processor.IExecutorProcessor[dto.MeRequest, dto.MeResponse]
	menusExecutor   processor.IExecutorProcessor[dto.MenusRequest, dto.MenusResponse]
	cookie          CookieConfig
}

type Executors struct {
	Login   processor.IExecutorProcessor[dto.LoginRequest, dto.LoginResult]
	Refresh processor.IExecutorProcessor[dto.RefreshRequest, dto.RefreshResult]
	Logout  processor.IExecutorProcessor[dto.LogoutRequest, processor.Empty]
	Me      processor.IExecutorProcessor[dto.MeRequest, dto.MeResponse]
	Menus   processor.IExecutorProcessor[dto.MenusRequest, dto.MenusResponse]
}

func NewAuthController(ex Executors, cookie CookieConfig) *AuthController {
	return &AuthController{
		loginExecutor:   ex.Login,
		refreshExecutor: ex.Refresh,
		logoutExecutor:  ex.Logout,
		meExecutor:      ex.Me,
		menusExecutor:   ex.Menus,
		cookie:          cookie,
	}
}

func (h *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, apperror.BadRequest("Data login tidak valid"))
		return
	}
	req.UserAgent = c.Request.UserAgent()
	req.IPAddress = c.ClientIP()

	result, err := h.loginExecutor.Process(c.Request.Context(), req)
	if err != nil {
		util.Error(c, err)
		return
	}

	h.setRefreshCookie(c, result.RefreshToken, result.RefreshExpiresAt)
	util.Success(c, http.StatusOK, result.Response)
}

func (h *AuthController) Refresh(c *gin.Context) {
	token, _ := c.Cookie(h.cookie.Name)

	result, err := h.refreshExecutor.Process(c.Request.Context(), dto.RefreshRequest{
		RefreshToken: token,
		UserAgent:    c.Request.UserAgent(),
		IPAddress:    c.ClientIP(),
	})
	if err != nil {
		util.Error(c, err)
		return
	}

	h.setRefreshCookie(c, result.RefreshToken, result.RefreshExpiresAt)
	util.Success(c, http.StatusOK, result.Response)
}

func (h *AuthController) Logout(c *gin.Context) {
	token, _ := c.Cookie(h.cookie.Name)

	if _, err := h.logoutExecutor.Process(c.Request.Context(), dto.LogoutRequest{RefreshToken: token}); err != nil {
		util.Error(c, err)
		return
	}

	h.clearRefreshCookie(c)
	util.Success(c, http.StatusOK, gin.H{})
}

func (h *AuthController) Me(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		util.Error(c, apperror.Unauthorized("Belum login"))
		return
	}

	result, err := h.meExecutor.Process(c.Request.Context(), dto.MeRequest{UserID: userID})
	if err != nil {
		util.Error(c, err)
		return
	}
	util.Success(c, http.StatusOK, result)
}

func (h *AuthController) Menus(c *gin.Context) {
	roleKode, ok := roleFromContext(c)
	if !ok {
		util.Error(c, apperror.Unauthorized("Belum login"))
		return
	}

	result, err := h.menusExecutor.Process(c.Request.Context(), dto.MenusRequest{RoleKode: roleKode})
	if err != nil {
		util.Error(c, err)
		return
	}
	util.Success(c, http.StatusOK, result)
}

func (h *AuthController) setRefreshCookie(c *gin.Context, token string, expiresAt time.Time) {
	maxAge := int(time.Until(expiresAt).Seconds())
	c.SetSameSite(parseSameSite(h.cookie.SameSite))
	c.SetCookie(h.cookie.Name, token, maxAge, refreshCookiePath, h.cookie.Domain, h.cookie.Secure, true)
}

func (h *AuthController) clearRefreshCookie(c *gin.Context) {
	c.SetSameSite(parseSameSite(h.cookie.SameSite))
	c.SetCookie(h.cookie.Name, "", -1, refreshCookiePath, h.cookie.Domain, h.cookie.Secure, true)
}

func parseSameSite(v string) http.SameSite {
	switch strings.ToLower(v) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}

func userIDFromContext(c *gin.Context) (uint64, bool) {
	v, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}
	id, ok := v.(uint64)
	return id, ok
}

func roleFromContext(c *gin.Context) (string, bool) {
	v, ok := c.Get("role")
	if !ok {
		return "", false
	}
	role, ok := v.(string)
	return role, ok && role != ""
}
