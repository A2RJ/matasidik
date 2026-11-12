package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/httpapi/middleware"
	"matasidik/backend/internal/httpapi/response"
	"matasidik/backend/internal/service"
)

const refreshCookie = "matasidik_refresh"

func (h *Handler) setRefreshCookie(c *gin.Context, token string, expires string) {
	exp, _ := time.Parse(time.RFC3339, expires)
	maxAge := int(time.Until(exp).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshCookie, token, maxAge, "/api/v1/auth", "", h.Cfg.CookieSecure, true)
}

func (h *Handler) clearRefreshCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshCookie, "", -1, "/api/v1/auth", "", h.Cfg.CookieSecure, true)
}

func (h *Handler) Login(c *gin.Context) {
	var in struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if !h.bind(c, &in) {
		return
	}
	result, err := h.Auth.Login(c.Request.Context(), in.Username, in.Password, c.Request.UserAgent())
	if err != nil {
		response.Error(c, err)
		return
	}
	h.setRefreshCookie(c, result.RefreshToken, result.RefreshExpires)
	response.OK(c, gin.H{
		"access_token": result.AccessToken, "access_expires_at": result.AccessExpires,
		"user": result.User,
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	refresh, err := c.Cookie(refreshCookie)
	if err != nil {
		response.Error(c, domain.Unauthorized("refresh token tidak ada"))
		return
	}
	result, err := h.Auth.Refresh(c.Request.Context(), refresh, c.Request.UserAgent())
	if err != nil {
		h.clearRefreshCookie(c)
		response.Error(c, err)
		return
	}
	h.setRefreshCookie(c, result.RefreshToken, result.RefreshExpires)
	response.OK(c, gin.H{
		"access_token": result.AccessToken, "access_expires_at": result.AccessExpires,
		"user": result.User,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	refresh, _ := c.Cookie(refreshCookie)
	_ = h.Auth.Logout(c.Request.Context(), refresh)
	h.clearRefreshCookie(c)
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Me(c *gin.Context) {
	user := h.actor(c)
	impersonator := middleware.ImpersonatorFrom(c)
	response.OK(c, gin.H{"user": user, "acting_as": impersonator})
}

func (h *Handler) Profile(c *gin.Context) {
	response.OK(c, gin.H{
		"user": h.actor(c), "acting_as": middleware.ImpersonatorFrom(c),
	})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	if middleware.ImpersonatorFrom(c) != nil {
		response.Error(c, domain.Forbidden("kembali ke akun admin untuk mengubah profil"))
		return
	}
	var in service.UpdateProfileInput
	if !h.bind(c, &in) {
		return
	}
	before := *h.actor(c)
	user, err := h.Auth.UpdateProfile(c.Request.Context(), h.actor(c), in)
	if err != nil {
		response.Error(c, err)
		return
	}
	_ = h.Audit.WriteWithImpersonator(c.Request.Context(), &user.ID, nil,
		"profile.updated", "user", user.ID, before, user, "", h.requestID(c))
	response.OK(c, user)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	if middleware.ImpersonatorFrom(c) != nil {
		response.Error(c, domain.Forbidden("kembali ke akun admin untuk mengubah password"))
		return
	}
	var in service.ChangePasswordInput
	if !h.bind(c, &in) {
		return
	}
	if err := h.Auth.ChangePassword(c.Request.Context(), h.actor(c), in); err != nil {
		response.Error(c, err)
		return
	}
	h.clearRefreshCookie(c)
	_ = h.Audit.WriteWithImpersonator(c.Request.Context(), &h.actor(c).ID, nil,
		"profile.password_changed", "user", h.actor(c).ID, nil, nil, "password diubah oleh pemilik akun", h.requestID(c))
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ActingAs(c *gin.Context) {
	var in struct {
		UserID string `json:"user_id"`
	}
	if !h.bind(c, &in) {
		return
	}
	result, err := h.Auth.ActingAs(c.Request.Context(), h.actor(c), in.UserID)
	if err != nil {
		response.Error(c, err)
		return
	}
	_ = h.Audit.WriteWithImpersonator(c.Request.Context(), &h.actor(c).ID, nil,
		"auth.acting_as", "user", in.UserID, nil, nil, "actingAs oleh admin", h.requestID(c))
	response.OK(c, gin.H{
		"access_token": result.AccessToken, "access_expires_at": result.AccessExpires,
		"user": result.User, "acting_as": result.ActingAs,
	})
}
