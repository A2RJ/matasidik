package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"matasidik/backend/internal/config"
	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/httpapi/middleware"
	"matasidik/backend/internal/repository"
	"matasidik/backend/internal/service"
)

type Handler struct {
	Cfg       config.Config
	Repo      repository.Repository
	Auth      *service.AuthService
	Admin     *service.AdminService
	Master    *service.MasterService
	Subs      *service.SubmissionService
	Calc      *service.CalculationService
	Dashboard *service.DashboardService
	Factor    *service.FactorService
	Audit     *service.AuditService
}

func New(cfg config.Config, repo repository.Repository, auth *service.AuthService, admin *service.AdminService, master *service.MasterService, subs *service.SubmissionService, calc *service.CalculationService, dashboard *service.DashboardService, factor *service.FactorService, audit *service.AuditService) *Handler {
	return &Handler{Cfg: cfg, Repo: repo, Auth: auth, Admin: admin, Master: master, Subs: subs, Calc: calc, Dashboard: dashboard, Factor: factor, Audit: audit}
}

func (h *Handler) actor(c *gin.Context) *domain.User { return middleware.ActorFrom(c) }

func (h *Handler) requestID(c *gin.Context) string { return middleware.GetRequestID(c) }

func (h *Handler) bind(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "VALIDATION_ERROR", "message": "body JSON tidak valid",
			"details": map[string]any{"error": err.Error()}, "request_id": h.requestID(c),
		})
		return false
	}
	return true
}

func limitFromQuery(c *gin.Context, fallback int) int {
	v, err := strconv.Atoi(c.Query("limit"))
	if err != nil || v <= 0 || v > 500 {
		return fallback
	}
	return v
}
