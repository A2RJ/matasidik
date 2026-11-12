package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRecoveryReturnsJSONOnPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID(), Recovery())
	router.GET("/panic", func(*gin.Context) {
		panic("uji panic")
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	if got := recorder.Header().Get("X-Request-ID"); got == "" {
		t.Fatal("X-Request-ID kosong")
	}
	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("body bukan JSON: %v (%q)", err, recorder.Body.String())
	}
	if body["code"] != "INTERNAL" {
		t.Fatalf("code = %v, want INTERNAL", body["code"])
	}
	if body["message"] == "" || body["request_id"] == "" {
		t.Fatalf("body tidak lengkap: %+v", body)
	}
}
