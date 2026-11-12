package handler

import (
	"github.com/gin-gonic/gin"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/httpapi/response"
)

func dashboardFilter(c *gin.Context) domain.DashboardFilter {
	return domain.DashboardFilter{
		PeriodeID: c.Query("periode"), ScopeType: c.Query("scope_type"), ScopeID: c.Query("scope_id"),
		FaskesID: c.Query("faskes_id"), BidangID: c.Query("bidang_id"),
		IndikatorID: c.Query("indikator_id"), Status: c.Query("status"),
	}
}

func (h *Handler) scopedDashboardFilter(c *gin.Context) domain.DashboardFilter {
	f := dashboardFilter(c)
	actor := h.actor(c)
	if actor != nil && (actor.Role == domain.RolePKMOperator || actor.Role == domain.RolePKMPJ) && actor.FaskesID != nil {
		f.FaskesID = *actor.FaskesID
		f.ScopeType = "FASKES"
		f.ScopeID = *actor.FaskesID
	}
	return f
}

func (h *Handler) DashboardSPM(c *gin.Context) {
	items, summary, err := h.Dashboard.SPM(c.Request.Context(), h.scopedDashboardFilter(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{
		"items": items, "summary": summary,
		"meta": gin.H{"status_basis": "target"},
	})
}

func (h *Handler) DashboardTrace(c *gin.Context) {
	trace, err := h.Dashboard.Trace(c.Request.Context(), c.Param("kode"), h.scopedDashboardFilter(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, trace)
}

func (h *Handler) IKUValues(c *gin.Context) {
	values, err := h.Dashboard.IKUValues(c.Request.Context(), c.Query("periode"), c.Query("scope_type"), c.Query("scope_id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, values)
}

func (h *Handler) IKUValueByKode(c *gin.Context) {
	value, err := h.Calc.ValueByKode(c.Request.Context(), c.Param("kode"), c.Query("periode"), c.DefaultQuery("scope_type", "KABUPATEN"), c.DefaultQuery("scope_id", "wil-kab-swb"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, value)
}

func (h *Handler) IKUValueDetail(c *gin.Context) {
	details, err := h.Calc.Detail(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, details)
}

func (h *Handler) DashboardDetail(c *gin.Context) {
	detail, err := h.Dashboard.Detail(c.Request.Context(), c.Param("kode"), h.scopedDashboardFilter(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, detail)
}

func (h *Handler) DashboardIKUChart(c *gin.Context) {
	chart, err := h.Dashboard.IKUByFaskes(c.Request.Context(), h.scopedDashboardFilter(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, chart)
}

func (h *Handler) ExportSPM(c *gin.Context) {
	actor := h.actor(c)
	if actor.Role != domain.RolePimpinan && actor.Role != domain.RoleKabid && actor.Role != domain.RoleSekretaris && actor.Role != domain.RoleVerifier && actor.Role != domain.RoleAdmin {
		response.Error(c, domain.Forbidden("role tidak boleh export dashboard"))
		return
	}
	file, err := h.Dashboard.ExportSPM(c.Request.Context(), h.scopedDashboardFilter(c), c.DefaultQuery("format", "xlsx"))
	if err != nil {
		response.Error(c, err)
		return
	}
	c.Header("Content-Type", file.ContentType)
	c.Header("Content-Disposition", `attachment; filename="`+file.Filename+`"`)
	c.Data(200, file.ContentType, file.Data)
}
