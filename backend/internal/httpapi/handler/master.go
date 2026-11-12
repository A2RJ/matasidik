package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"matasidik/backend/internal/httpapi/response"
	"matasidik/backend/internal/service"
)

func (h *Handler) Health(c *gin.Context) {
	response.OK(c, gin.H{"status": "ok", "service": "MATA SIDIK API"})
}

func (h *Handler) Bootstrap(c *gin.Context) {
	data, err := h.Master.Bootstrap(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) Actors(c *gin.Context) {
	users, err := h.Repo.ListUsers(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, users)
}

func (h *Handler) Labels(c *gin.Context) {
	data, err := h.Master.Labels(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) CreateLabel(c *gin.Context) {
	var in service.CreateLabelInput
	if !h.bind(c, &in) {
		return
	}
	label, err := h.Master.CreateLabel(c.Request.Context(), in, h.actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	_ = h.Audit.Write(c.Request.Context(), &h.actor(c).ID, "master.label_created", "label", label.ID, nil, label, "", h.requestID(c))
	response.Created(c, label)
}

func (h *Handler) Indicators(c *gin.Context) {
	data, err := h.Master.Indicators(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) CreateIndicator(c *gin.Context) {
	var in service.CreateIndicatorInput
	if !h.bind(c, &in) {
		return
	}
	ind, err := h.Master.CreateIndicator(c.Request.Context(), in, h.actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	_ = h.Audit.Write(c.Request.Context(), &h.actor(c).ID, "master.indicator_created", "indikator", ind.ID, nil, ind, "", h.requestID(c))
	response.Created(c, ind)
}

func (h *Handler) Targets(c *gin.Context) {
	data, err := h.Master.Targets(c.Request.Context(), c.Query("periode"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) UpsertTarget(c *gin.Context) {
	var in service.UpsertTargetInput
	if !h.bind(c, &in) {
		return
	}
	target, err := h.Master.UpsertTarget(c.Request.Context(), in, h.actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	_ = h.Audit.Write(c.Request.Context(), &h.actor(c).ID, "master.target_upserted", "target", target.ID, nil, target, "", h.requestID(c))
	response.OK(c, target)
}

func (h *Handler) Periods(c *gin.Context) {
	data, err := h.Master.Periods(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) CapaianTahunanList(c *gin.Context) {
	tahun, _ := strconv.Atoi(c.Query("tahun"))
	if tahun == 0 {
		tahun = time.Now().Year() - 1
	}
	data, err := h.Master.CapaianTahunan(c.Request.Context(), tahun)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) CreateCapaianTahunan(c *gin.Context) {
	var in service.CapaianTahunanInput
	if !h.bind(c, &in) {
		return
	}
	data, err := h.Master.CreateCapaianTahunan(c.Request.Context(), in, h.actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	_ = h.Audit.WriteWithImpersonator(c.Request.Context(), &h.actor(c).ID, nil,
		"master.capaian_tahunan_created", "capaian_tahunan", data.ID, nil, data, "", h.requestID(c))
	response.Created(c, data)
}

func (h *Handler) UpdateCapaianTahunan(c *gin.Context) {
	var in service.CapaianTahunanInput
	if !h.bind(c, &in) {
		return
	}
	before, _ := h.Repo.GetCapaianTahunan(c.Request.Context(), c.Param("id"))
	data, err := h.Master.UpdateCapaianTahunan(c.Request.Context(), c.Param("id"), in, h.actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	_ = h.Audit.WriteWithImpersonator(c.Request.Context(), &h.actor(c).ID, nil,
		"master.capaian_tahunan_updated", "capaian_tahunan", data.ID, before, data, "", h.requestID(c))
	response.OK(c, data)
}

func (h *Handler) DeleteCapaianTahunan(c *gin.Context) {
	before, _ := h.Repo.GetCapaianTahunan(c.Request.Context(), c.Param("id"))
	if err := h.Master.DeleteCapaianTahunan(c.Request.Context(), c.Param("id"), h.actor(c)); err != nil {
		response.Error(c, err)
		return
	}
	_ = h.Audit.WriteWithImpersonator(c.Request.Context(), &h.actor(c).ID, nil,
		"master.capaian_tahunan_deleted", "capaian_tahunan", c.Param("id"), before, nil, "", h.requestID(c))
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Template(c *gin.Context) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `attachment; filename="template-bnba.csv"`)
	c.String(http.StatusOK, "periode;kode_faskes;kode_iku;kode_indikator;scope_type;scope_code;nilai;numerator;denominator;catatan\n")
}
