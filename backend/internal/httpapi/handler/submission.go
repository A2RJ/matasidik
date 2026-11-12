package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/httpapi/response"
	"matasidik/backend/internal/service"
)

func (h *Handler) UploadSubmission(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.Cfg.MaxUploadBytes()+(1<<20))
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, domain.Validation("file wajib diunggah", map[string]any{"error": err.Error()}))
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, err)
		return
	}
	defer file.Close()
	sub, job, err := h.Subs.Upload(c.Request.Context(), h.actor(c), service.UploadInput{
		File: file, Filename: fileHeader.Filename, Size: fileHeader.Size,
		PeriodeKode: c.PostForm("periode"), FaskesID: c.PostForm("faskes_id"),
		IdempotencyKey: c.GetHeader("Idempotency-Key"), RequestID: h.requestID(c),
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	if job == nil {
		response.OK(c, gin.H{"submission": sub})
		return
	}
	response.Accepted(c, gin.H{"submission": sub, "job": job})
}

func (h *Handler) Job(c *gin.Context) {
	job, err := h.Repo.GetJob(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	actor := h.actor(c)
	if actor.Role == domain.RolePKMOperator || actor.Role == domain.RolePKMPJ {
		var payload struct {
			SubmissionID string `json:"submission_id"`
		}
		if err := json.Unmarshal([]byte(job.PayloadJSON), &payload); err != nil {
			response.Error(c, domain.Forbidden("job di luar scope"))
			return
		}
		sub, err := h.Repo.GetSubmission(c.Request.Context(), payload.SubmissionID)
		if err != nil || actor.FaskesID == nil || sub.FaskesID != *actor.FaskesID {
			response.Error(c, domain.Forbidden("job di luar scope faskes"))
			return
		}
	}
	response.OK(c, job)
}

func (h *Handler) Submissions(c *gin.Context) {
	filter := domain.SubmissionFilter{
		Status: c.Query("status"), PeriodeID: c.Query("periode"),
		FaskesID: c.Query("faskes_id"), Limit: limitFromQuery(c, 100),
	}
	actor := h.actor(c)
	if actor.Role == domain.RolePKMOperator || actor.Role == domain.RolePKMPJ {
		if actor.FaskesID == nil {
			response.OK(c, []domain.Submission{})
			return
		}
		filter.FaskesID = *actor.FaskesID
	}
	list, err := h.Subs.List(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, list)
}

func (h *Handler) SubmissionDetail(c *gin.Context) {
	detail, err := h.Subs.Detail(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	actor := h.actor(c)
	if actor.Role == domain.RolePKMOperator || actor.Role == domain.RolePKMPJ {
		if actor.FaskesID == nil || detail.Submission.FaskesID != *actor.FaskesID {
			response.Error(c, domain.Forbidden("submission di luar scope faskes"))
			return
		}
	}
	response.OK(c, detail)
}

func (h *Handler) UpdateSubmissionRow(c *gin.Context) {
	var in service.UpdateRowInput
	if !h.bind(c, &in) {
		return
	}
	row, err := h.Subs.UpdateRow(c.Request.Context(), h.actor(c), c.Param("id"), c.Param("rowId"), in, h.requestID(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) SubmitSubmission(c *gin.Context) {
	sub, err := h.Subs.Submit(c.Request.Context(), h.actor(c), c.Param("id"), h.requestID(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, sub)
}

func (h *Handler) Decision(c *gin.Context) {
	var in service.DecisionInput
	if !h.bind(c, &in) {
		return
	}
	in.RequestID = h.requestID(c)
	sub, err := h.Subs.Decision(c.Request.Context(), h.actor(c), c.Param("id"), in)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, sub)
}

func (h *Handler) Revisions(c *gin.Context) {
	data, err := h.Subs.Revisions(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) File(c *gin.Context) {
	file, err := h.Repo.GetFileObject(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	actor := h.actor(c)
	if actor.Role == domain.RolePKMOperator || actor.Role == domain.RolePKMPJ {
		sub, err := h.Repo.GetSubmissionByFileID(c.Request.Context(), file.ID)
		if err != nil || actor.FaskesID == nil || sub.FaskesID != *actor.FaskesID {
			response.Error(c, domain.Forbidden("file di luar scope faskes"))
			return
		}
	}
	name := file.OriginalName
	if name == "" {
		name = "download"
	}
	c.Header("Content-Disposition", `attachment; filename="`+strings.ReplaceAll(name, `"`, "")+`"`)
	if file.Mime != nil {
		c.Header("Content-Type", *file.Mime)
	}
	c.File(file.Path)
}
