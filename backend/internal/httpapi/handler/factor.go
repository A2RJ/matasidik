package handler

import (
	"github.com/gin-gonic/gin"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/httpapi/response"
	"matasidik/backend/internal/service"
)

func (h *Handler) Factors(c *gin.Context) {
	rows, err := h.Factor.List(c.Request.Context(), c.Query("periode"), c.DefaultQuery("scope_type", "KABUPATEN"), c.DefaultQuery("scope_id", "wil-kab-swb"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, rows)
}

func (h *Handler) UpsertFactor(c *gin.Context) {
	var in service.UpsertFactorInput
	if !h.bind(c, &in) {
		return
	}
	in.IKUValueID = c.Param("ikuValueID")
	in.RequestID = h.requestID(c)
	f, err := h.Factor.Upsert(c.Request.Context(), h.actor(c), in)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, f)
}

func (h *Handler) TindakLanjut(c *gin.Context) {
	rows, err := h.Factor.ListTL(c.Request.Context(), c.Query("periode"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, rows)
}

func (h *Handler) CreateTindakLanjut(c *gin.Context) {
	var in service.TindakLanjutInput
	if !h.bind(c, &in) {
		return
	}
	in.RequestID = h.requestID(c)
	tl, err := h.Factor.CreateTL(c.Request.Context(), h.actor(c), in)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, tl)
}

func (h *Handler) UpdateTindakLanjut(c *gin.Context) {
	var in service.TindakLanjutInput
	if !h.bind(c, &in) {
		return
	}
	in.RequestID = h.requestID(c)
	tl, err := h.Factor.UpdateTL(c.Request.Context(), h.actor(c), c.Param("id"), in)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, tl)
}

func (h *Handler) AuditLogs(c *gin.Context) {
	actor := h.actor(c)
	if actor.Role != domain.RoleAdmin && actor.Role != domain.RolePimpinan {
		response.Error(c, domain.Forbidden("hanya admin/pimpinan yang boleh melihat audit"))
		return
	}
	logs, err := h.Audit.List(c.Request.Context(), domain.AuditFilter{
		EntityType: c.Query("entity_type"), EntityID: c.Query("entity_id"),
		Limit: limitFromQuery(c, 200),
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, logs)
}
