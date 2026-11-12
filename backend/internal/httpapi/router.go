package httpapi

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"matasidik/backend/internal/config"
	"matasidik/backend/internal/httpapi/handler"
	"matasidik/backend/internal/httpapi/middleware"
)

func NewRouter(cfg config.Config, h *handler.Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.MaxMultipartMemory = 16 << 20
	r.Use(gin.Logger(), middleware.RequestID(), middleware.Recovery(), middleware.CORS())

	r.GET("/healthz", h.Health)
	r.GET("/readyz", h.Health)

	publicAuth := r.Group("/api/v1/auth")
	{
		publicAuth.POST("/login", h.Login)
		publicAuth.POST("/refresh", h.Refresh)
	}

	api := r.Group("/api/v1")
	api.Use(middleware.AuthRequired(h.Auth))
	{
		api.POST("/auth/logout", h.Logout)
		api.GET("/auth/me", h.Me)
		api.POST("/auth/acting-as", h.ActingAs)
		api.GET("/profile", h.Profile)
		api.PATCH("/profile", h.UpdateProfile)
		api.POST("/profile/change-password", h.ChangePassword)
		api.GET("/meta/bootstrap", h.Bootstrap)
		api.GET("/master/labels", h.Labels)
		api.POST("/master/labels", h.CreateLabel)
		api.GET("/master/indicators", h.Indicators)
		api.POST("/master/indicators", h.CreateIndicator)
		api.GET("/master/targets", h.Targets)
		api.PUT("/master/targets", h.UpsertTarget)
		api.GET("/master/periods", h.Periods)
		api.GET("/master/template", h.Template)
		api.GET("/master/capaian-tahunan", h.CapaianTahunanList)
		api.POST("/master/capaian-tahunan", h.CreateCapaianTahunan)
		api.PATCH("/master/capaian-tahunan/:id", h.UpdateCapaianTahunan)
		api.DELETE("/master/capaian-tahunan/:id", h.DeleteCapaianTahunan)

		api.POST("/submissions", h.UploadSubmission)
		api.GET("/submissions", h.Submissions)
		api.GET("/submissions/:id", h.SubmissionDetail)
		api.PATCH("/submissions/:id/rows/:rowId", h.UpdateSubmissionRow)
		api.POST("/submissions/:id/submit", h.SubmitSubmission)
		api.POST("/submissions/:id/decision", h.Decision)
		api.GET("/submissions/:id/revisions", h.Revisions)
		api.GET("/jobs/:id", h.Job)

		api.GET("/dashboard/spm", h.DashboardSPM)
		api.GET("/dashboard/spm/export", h.ExportSPM)
		api.GET("/dashboard/spm/:kode/detail", h.DashboardDetail)
		api.GET("/dashboard/spm/:kode/trace", h.DashboardTrace)
		api.GET("/dashboard/iku-by-faskes", h.DashboardIKUChart)
		api.GET("/iku/values", h.IKUValues)
		api.GET("/iku/:kode/values", h.IKUValueByKode)
		api.GET("/iku-values/:id/detail", h.IKUValueDetail)

		api.GET("/factors", h.Factors)
		api.PUT("/factors/:ikuValueID", h.UpsertFactor)
		api.GET("/tindak-lanjut", h.TindakLanjut)
		api.POST("/tindak-lanjut", h.CreateTindakLanjut)
		api.PATCH("/tindak-lanjut/:id", h.UpdateTindakLanjut)

		api.GET("/audit", h.AuditLogs)
		api.GET("/files/:id", h.File)

		api.GET("/admin/overview", h.AdminOverview)
		api.GET("/admin/users", h.AdminUsers)
		api.POST("/admin/users", h.AdminCreateUser)
		api.PATCH("/admin/users/:id", h.AdminUpdateUser)
		api.POST("/admin/users/:id/reset-password", h.AdminResetPassword)
		api.GET("/admin/faskes", h.AdminFaskes)
		api.POST("/admin/faskes", h.AdminCreateFaskes)
		api.PATCH("/admin/faskes/:id", h.AdminUpdateFaskes)
		api.GET("/admin/positions", h.AdminPositions)
		api.POST("/admin/positions", h.AdminCreatePosition)
		api.PATCH("/admin/positions/:id", h.AdminUpdatePosition)
	}

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "endpoint tidak ditemukan", "request_id": middleware.GetRequestID(c)})
			return
		}
		webDir := cfg.WebDir
		if webDir == "" {
			c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "frontend build tidak ditemukan"})
			return
		}
		clean := filepath.Clean(c.Request.URL.Path)
		target := filepath.Join(webDir, clean)
		if info, err := os.Stat(target); err == nil && !info.IsDir() {
			c.File(target)
			return
		}
		index := filepath.Join(webDir, "index.html")
		if _, err := os.Stat(index); err == nil {
			c.File(index)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "halaman tidak ditemukan"})
	})
	return r
}
