package handler

import (
	"github.com/gin-gonic/gin"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/httpapi/response"
	"matasidik/backend/internal/service"
)

func (h *Handler) CreateIKUChange(c *gin.Context) {
	var in service.IKUChangeInput
	if !h.bind(c, &in) {
		return
	}
	in.RequestID = h.requestID(c)
	data, err := h.IKU.Create(c.Request.Context(), h.actor(c), in)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, data)
}

func (h *Handler) IKUChanges(c *gin.Context) {
	filter := domain.IKUChangeFilter{
		Status: c.Query("status"), FaskesID: c.Query("faskes_id"),
		RequestedBy: c.Query("requested_by"), Limit: limitFromQuery(c, 100),
	}
	data, err := h.IKU.List(c.Request.Context(), h.actor(c), filter)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) IKUChangeDetail(c *gin.Context) {
	request, revisions, reviews, err := h.IKU.Detail(c.Request.Context(), h.actor(c), c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"request": request, "revisions": revisions, "reviews": reviews})
}

func (h *Handler) DecideIKUChange(c *gin.Context) {
	var in service.IKUChangeDecisionInput
	if !h.bind(c, &in) {
		return
	}
	in.RequestID = h.requestID(c)
	data, err := h.IKU.Decide(c.Request.Context(), h.actor(c), c.Param("id"), in)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, data)
}
