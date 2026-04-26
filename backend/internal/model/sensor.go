package model

type SensorData struct {
	Distance float64 `json:"distance"`
	Temperature float64 `json:"temperature"`
	Pressure float64 `json:"pressure"`
}