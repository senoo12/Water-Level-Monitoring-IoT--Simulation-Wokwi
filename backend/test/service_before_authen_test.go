package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_BeforeAuthentication(t *testing.T) {

	gin.SetMode(gin.TestMode)

	r := gin.Default()

	// endpoint TANPA middleware auth
	r.POST("/sensor", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
		})
	})

	payload := []byte(`{
		"distance": 20,
		"temperature": 29,
		"pressure": 1001
	}`)

	req, _ := http.NewRequest(
		"POST",
		"/sensor",
		bytes.NewBuffer(payload),
	)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}