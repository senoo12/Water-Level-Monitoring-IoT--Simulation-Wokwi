package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"water-monitor/internal/config"
	"water-monitor/internal/handler"
	"water-monitor/internal/repository"
	"water-monitor/internal/service"
)

func main() {
	godotenv.Load()
	cfg := config.LoadConfig()

	repo := repository.NewThingSpeakRepository(cfg.ThingSpeakAPIKey)
	service := service.NewSensorService(repo)
	handler := handler.NewSensorHandler(service)

	r := gin.Default()

	r.POST("/sensor", handler.ReceiveData)

	r.Run(":8080")
}