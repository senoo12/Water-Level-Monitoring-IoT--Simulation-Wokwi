package handler

import (
	"github.com/gin-gonic/gin"
	"water-monitor/internal/model"
	"water-monitor/internal/service"
	"water-monitor/pkg/response"
	"encoding/json"
)

type SensorHandler struct {
	service *service.SensorService
}

func NewSensorHandler(service *service.SensorService) *SensorHandler {
	return &SensorHandler{service: service}
}

func (h *SensorHandler) ReceiveData(c *gin.Context) {
	var data model.SensorData

	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&data); err != nil {
		response.Error(c, 400, "invalid request body: "+err.Error())
		return
	}

	status, statusCode, waterLevel, err := h.service.ProcessData(data)
	if err != nil {
		response.Error(c, 500, "failed to send data")
		return
	}

	if data.Distance < 0 || data.Distance > 100 {
		response.Error(c, 400, "invalid distance value")
		return
	}

	response.Success(c, "success", gin.H{
		"distance":    data.Distance,
		"water_level": waterLevel,
		"temperature": data.Temperature,
		"pressure":    data.Pressure,
		"status":      status,
		"status_code": statusCode,
	})
}