package service

import (
	"fmt"
	"water-monitor/internal/model"
	"water-monitor/internal/repository"
)

type SensorService struct {
	repo repository.SensorRepository
}

func NewSensorService(repo repository.SensorRepository) *SensorService {
	return &SensorService{repo: repo}
}

func (s *SensorService) ProcessData(data model.SensorData) (string, int, float64, error) {
	const tankHeight = 100.0

	// ===== VALIDASI INPUT =====
	if data.Distance < 0 || data.Distance > 1000 {
		return "", 0, 0, fmt.Errorf("invalid distance value")
	}

	// ===== HITUNG WATER LEVEL =====
	waterLevel := tankHeight - data.Distance

	// clamp
	if waterLevel < 0 {
		waterLevel = 0
	}
	if waterLevel > tankHeight {
		waterLevel = tankHeight
	}

	// ===== FUZZY MEMBERSHIP =====
	muAman := membershipAman(waterLevel)
	muWaspada := membershipWaspada(waterLevel)
	muBahaya := membershipBahaya(waterLevel)

	// ===== DEBUG (opsional) =====
	fmt.Printf("WaterLevel: %.2f | Aman: %.2f | Waspada: %.2f | Bahaya: %.2f\n",
		waterLevel, muAman, muWaspada, muBahaya)

	// ===== DEFUZZIFICATION (MAX MEMBERSHIP) =====
	status := "AMAN"
	statusCode := 1
	maxVal := muAman

	if muWaspada > maxVal {
		status = "WASPADA"
		statusCode = 2
		maxVal = muWaspada
	}

	if muBahaya > maxVal {
		status = "BAHAYA"
		statusCode = 3
	}

	// ===== KIRIM KE THINGSPEAK =====
	err := s.repo.SendData(
		waterLevel,
		data.Temperature,
		data.Pressure,
		data.Distance,
		statusCode,
	)
	if err != nil {
		fmt.Println("ERROR SEND DATA:", err)
		return "", 0, 0, err
	}

	// ===== SIMPAN KE GOOGLE SHEET =====
	err = repository.SaveToGoogleSheet(
		waterLevel,
		data.Temperature,
		data.Pressure,
		data.Distance,
		status,
	)

	if err != nil {
		fmt.Println("ERROR SEND DATA:", err)
		return "", 0, 0, err
	}

	return status, statusCode, waterLevel, nil
}

//
// ===== FUZZY MEMBERSHIP FUNCTIONS =====
//

// AMAN (0 - 30)
func membershipAman(x float64) float64 {
	if x <= 0 {
		return 1
	}
	if x > 0 && x < 30 {
		return (30 - x) / 30
	}
	return 0
}

// WASPADA (30 - 70, peak 50)
func membershipWaspada(x float64) float64 {
	if x <= 30 || x >= 70 {
		return 0
	}
	if x > 30 && x <= 50 {
		return (x - 30) / 20
	}
	if x > 50 && x < 70 {
		return (70 - x) / 20
	}
	return 0
}

// BAHAYA (70 - 100)
func membershipBahaya(x float64) float64 {
	if x <= 70 {
		return 0
	}
	if x > 70 && x < 100 {
		return (x - 70) / 30
	}
	return 1
}
