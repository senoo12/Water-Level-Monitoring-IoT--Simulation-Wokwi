package main

import (
	"water-monitor/internal/config"
	"water-monitor/internal/handler"
	"water-monitor/internal/middleware"
	"water-monitor/internal/repository"
	"water-monitor/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg := config.LoadConfig()

	repo := repository.NewThingSpeakRepository(
		cfg.ThingSpeakAPIKey,
	)

	sensorService := service.NewSensorService(repo)

	sensorHandler := handler.NewSensorHandler(
		sensorService,
	)

	r := gin.Default()

	// route tanpa auth
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "server running",
		})
	})

	// route dengan auth middleware
	r.POST(
		"/sensor",
		middleware.AuthMiddleware(),
		sensorHandler.ReceiveData,
	)

	r.Run(":8080")
}