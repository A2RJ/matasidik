package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/httpapi/middleware"
)

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data, "request_id": middleware.GetRequestID(c)})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{"data": data, "request_id": middleware.GetRequestID(c)})
}

func Accepted(c *gin.Context, data any) {
	c.JSON(http.StatusAccepted, gin.H{"data": data, "request_id": middleware.GetRequestID(c)})
}

func NoContent(c *gin.Context) { c.Status(http.StatusNoContent) }

func Error(c *gin.Context, err error) {
	var app *domain.AppError
	if errors.As(err, &app) {
		status := app.Status
		if status == 0 {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{
			"code": app.Code, "message": app.Message, "details": app.Details,
			"request_id": middleware.GetRequestID(c),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"code": "INTERNAL", "message": "terjadi kesalahan internal",
		"details": map[string]any{"cause": err.Error()}, "request_id": middleware.GetRequestID(c),
	})
}
