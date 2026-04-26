package repository

import (
	"fmt"
	"io"
	"net/http"
)

type ThingSpeakRepository struct {
	APIKey string
}

func NewThingSpeakRepository(apiKey string) *ThingSpeakRepository {
	return &ThingSpeakRepository{APIKey: apiKey}
}

func (r *ThingSpeakRepository) SendData(
	water, temp, pressure, distance float64,
	status int,
) error {

	// ===== VALIDASI API KEY =====
	if r.APIKey == "" {
		return fmt.Errorf("thingspeak api key is empty")
	}

	// ===== BUILD URL =====
	url := fmt.Sprintf(
		"https://api.thingspeak.com/update?api_key=%s&field1=%.2f&field2=%.2f&field3=%.2f&field4=%.2f&field5=%d",
		r.APIKey,
		water,
		temp,
		pressure,
		distance,
		status,
	)

	fmt.Println("DEBUG URL:", url)

	// ===== REQUEST =====
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// ===== READ RESPONSE =====
	body, _ := io.ReadAll(resp.Body)

	fmt.Println("ThingSpeak Status:", resp.Status)
	fmt.Println("ThingSpeak Body:", string(body))

	// ===== VALIDASI RESPONSE =====
	if string(body) == "0" {
		return fmt.Errorf("thingspeak rejected data (0 response)")
	}

	return nil
}