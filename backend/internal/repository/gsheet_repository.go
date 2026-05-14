package repository

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type GSheetPayload struct {
	WaterLevel float64 `json:"water_level"`
	Temperature float64 `json:"temperature"`
	Pressure float64 `json:"pressure"`
	Distance float64 `json:"distance"`
	Status string `json:"status"`
}

func SaveToGoogleSheet(
	water float64,
	temp float64,
	pressure float64,
	distance float64,
	status string,
) error {

	payload := GSheetPayload{
		WaterLevel: water,
		Temperature: temp,
		Pressure: pressure,
		Distance: distance,
		Status: status,
	}

	jsonData, _ := json.Marshal(payload)

	_, err := http.Post(
		os.Getenv("GSHEET_WEBHOOK"),
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	return err
}