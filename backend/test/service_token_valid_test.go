package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"water-monitor/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func Test_AfterAuthentication_WithValidToken(t *testing.T) {

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

	// ===== CREATE VALID JWT =====
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": 1,
			"exp": time.Now().Add(time.Hour).Unix(),
		},
	)

	tokenString, _ := token.SignedString(
		middleware.SECRET_KEY,
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
		"Bearer "+tokenString,
	)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}