package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"water-monitor/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_AfterAuthentication_InvalidToken(t *testing.T) {

	gin.SetMode(gin.TestMode)

	r := gin.Default()

	r.POST(
		"/sensor",
		middleware.AuthMiddleware(),
		func(c *gin.Context) {
			c.JSON(200, gin.H{
				"success": true,
			})
		},
	)

	payload := []byte(`{
		"distance": 20
	}`)

	req, _ := http.NewRequest(
		"POST",
		"/sensor",
		bytes.NewBuffer(payload),
	)

	req.Header.Set(
		"Authorization",
		"Bearer WRONG_TOKEN",
	)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}