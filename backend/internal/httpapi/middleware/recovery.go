package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery menangkap panic dan tetap membalas JSON agar klien tidak menerima
// response 500 kosong yang sulit ditampilkan ke pengguna.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				slog.Error("panic pada request",
					"error", recovered,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"request_id", GetRequestID(c),
					"stack", string(debug.Stack()),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":       "INTERNAL",
					"message":    "terjadi kesalahan pada server",
					"request_id": GetRequestID(c),
				})
			}
		}()
		c.Next()
	}
}
