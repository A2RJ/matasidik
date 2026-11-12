package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/service"
)

const (
	actorKey        = "actor"
	impersonatorKey = "impersonator"
)

func AuthRequired(auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": "UNAUTHENTICATED", "message": "token tidak ada", "request_id": GetRequestID(c),
			})
			return
		}
		raw := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		ac, err := auth.Parse(c.Request.Context(), raw)
		if err != nil {
			status := http.StatusUnauthorized
			if app, ok := err.(*domain.AppError); ok {
				status = app.Status
			}
			c.AbortWithStatusJSON(status, gin.H{
				"code": "UNAUTHENTICATED", "message": "sesi tidak valid", "request_id": GetRequestID(c),
			})
			return
		}
		c.Set(actorKey, ac.User)
		if ac.Impersonator != nil {
			c.Set(impersonatorKey, ac.Impersonator)
			c.Request = c.Request.WithContext(service.WithImpersonator(c.Request.Context(), ac.Impersonator))
		}
		c.Next()
	}
}

func ActorFrom(c *gin.Context) *domain.User {
	if v, ok := c.Get(actorKey); ok {
		if u, ok := v.(*domain.User); ok {
			return u
		}
	}
	return nil
}

func ImpersonatorFrom(c *gin.Context) *domain.User {
	if v, ok := c.Get(impersonatorKey); ok {
		if u, ok := v.(*domain.User); ok {
			return u
		}
	}
	return nil
}
