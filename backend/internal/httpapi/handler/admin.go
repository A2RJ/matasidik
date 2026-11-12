package handler

import (
	"github.com/gin-gonic/gin"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/httpapi/response"
	"matasidik/backend/internal/service"
)

func (h *Handler) requireAdmin(c *gin.Context) bool {
	if h.actor(c).Role != domain.RoleAdmin {
		response.Error(c, domain.Forbidden("hanya admin"))
		return false
	}
	return true
}

func (h *Handler) AdminOverview(c *gin.Context) {
	if !h.requireAdmin(c) {
		return
	}
	users, err := h.Admin.Users(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	faskes, err := h.Admin.Faskes(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	positions, err := h.Admin.Positions(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	bidang, err := h.Repo.ListBidang(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"users": users, "faskes": faskes, "positions": positions, "bidang": bidang})
}

func (h *Handler) AdminUsers(c *gin.Context) {
	if !h.requireAdmin(c) {
		return
	}
	users, err := h.Admin.Users(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, users)
}

func (h *Handler) AdminCreateUser(c *gin.Context) {
	if !h.requireAdmin(c) {
		return
	}
	var in service.CreateUserInput
	if !h.bind(c, &in) {
		return
	}
	user, err := h.Admin.CreateUser(c.Request.Context(), in, h.actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	_ = h.Audit.WriteWithImpersonator(c.Request.Context(), &h.actor(c).ID, nil,
		"admin.user_created", "user", user.ID, nil, user, "", h.requestID(c))
	response.Created(c, user)
}

func (h *Handler) AdminUpdateUser(c *gin.Context) {
	if !h.requireAdmin(c) {
		return
	}
	var in service.UpdateUserInput
	if !h.bind(c, &in) {
		return
	}
	user, err := h.Admin.UpdateUser(c.Request.Context(), c.Param("id"), in, h.actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	_ = h.Audit.WriteWithImpersonator(c.Request.Context(), &h.actor(c).ID, nil,
		"admin.user_updated", "user", user.ID, nil, user, "", h.requestID(c))
	response.OK(c, user)
}

func (h *Handler) AdminResetPassword(c *gin.Context) {
	if !h.requireAdmin(c) {
		return
	}
	var in struct {
		Password string `json:"password"`
	}
	if !h.bind(c, &in) {
		return
	}
	if err := h.Admin.ResetPassword(c.Request.Context(), c.Param("id"), in.Password, h.actor(c)); err != nil {
		response.Error(c, err)
		return
	}
	_ = h.Audit.WriteWithImpersonator(c.Request.Context(), &h.actor(c).ID, nil,
		"admin.user_password_reset", "user", c.Param("id"), nil, nil, "", h.requestID(c))
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) AdminFaskes(c *gin.Context) {
	if !h.requireAdmin(c) {
		return
	}
	faskes, err := h.Admin.Faskes(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, faskes)
}

func (h *Handler) AdminCreateFaskes(c *gin.Context) {
	if !h.requireAdmin(c) {
		return
	}
	var in service.FaskesInput
	if !h.bind(c, &in) {
		return
	}
	f, err := h.Admin.CreateFaskes(c.Request.Context(), in, h.actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	_ = h.Audit.WriteWithImpersonator(c.Request.Context(), &h.actor(c).ID, nil,
		"admin.faskes_created", "faskes", f.ID, nil, f, "", h.requestID(c))
	response.Created(c, f)
}

func (h *Handler) AdminUpdateFaskes(c *gin.Context) {
	if !h.requireAdmin(c) {
		return
	}
	var in service.FaskesInput
	if !h.bind(c, &in) {
		return
	}
	f, err := h.Admin.UpdateFaskes(c.Request.Context(), c.Param("id"), in, h.actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	_ = h.Audit.WriteWithImpersonator(c.Request.Context(), &h.actor(c).ID, nil,
		"admin.faskes_updated", "faskes", f.ID, nil, f, "", h.requestID(c))
	response.OK(c, f)
}

func (h *Handler) AdminPositions(c *gin.Context) {
	if !h.requireAdmin(c) {
		return
	}
	positions, err := h.Admin.Positions(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, positions)
}

func (h *Handler) AdminCreatePosition(c *gin.Context) {
	if !h.requireAdmin(c) {
		return
	}
	var in service.PositionInput
	if !h.bind(c, &in) {
		return
	}
	p, err := h.Admin.CreatePosition(c.Request.Context(), in, h.actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	_ = h.Audit.WriteWithImpersonator(c.Request.Context(), &h.actor(c).ID, nil,
		"admin.position_created", "position", p.ID, nil, p, "", h.requestID(c))
	response.Created(c, p)
}

func (h *Handler) AdminUpdatePosition(c *gin.Context) {
	if !h.requireAdmin(c) {
		return
	}
	var in service.PositionInput
	if !h.bind(c, &in) {
		return
	}
	p, err := h.Admin.UpdatePosition(c.Request.Context(), c.Param("id"), in, h.actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	_ = h.Audit.WriteWithImpersonator(c.Request.Context(), &h.actor(c).ID, nil,
		"admin.position_updated", "position", p.ID, nil, p, "", h.requestID(c))
	response.OK(c, p)
}
