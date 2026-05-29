package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"water-monitor/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_BruteForce_InvalidToken(t *testing.T) {

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

	for i := 0; i < 10; i++ {

		req, _ := http.NewRequest(
			"POST",
			"/sensor",
			nil,
		)

		req.Header.Set(
			"Authorization",
			"Bearer WRONG_TOKEN",
		)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
	}
}