package test

import (
	"testing"
	"water-monitor/internal/model"
	"water-monitor/internal/service"
)

// ===== MOCK REPOSITORY =====
type MockRepo struct{}

func (m *MockRepo) SendData(water, temp, pressure, distance float64, status int) error {
	return nil
}

func TestProcessData_FuzzyLogic(t *testing.T) {
	repo := &MockRepo{}
	svc := service.NewSensorService(repo)

	tests := []struct {
		name           string
		input          model.SensorData
		expectedStatus string
	}{
		{
			name: "AMAN - water rendah",
			input: model.SensorData{
				Distance: 90,
				Temperature: 28,
				Pressure: 1012,
			},
			expectedStatus: "AMAN",
		},
		{
			name: "WASPADA - water sedang",
			input: model.SensorData{
				Distance: 50,
				Temperature: 28,
				Pressure: 1012,
			},
			expectedStatus: "WASPADA",
		},
		{
			name: "BAHAYA - water tinggi",
			input: model.SensorData{
				Distance: 10,
				Temperature: 28,
				Pressure: 1012,
			},
			expectedStatus: "BAHAYA",
		},
		{
			name: "TRANSISI - mendekati bahaya",
			input: model.SensorData{
				Distance: 25,
				Temperature: 28,
				Pressure: 1012,
			},
			expectedStatus: "BAHAYA",
		},
		{
			name: "TRANSISI - mendekati aman",
			input: model.SensorData{
				Distance: 65,
				Temperature: 28,
				Pressure: 1012,
			},
			expectedStatus: "WASPADA",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			status, code, waterLevel, err := svc.ProcessData(tt.input)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if status != tt.expectedStatus {
				t.Errorf("expected status %s, got %s", tt.expectedStatus, status)
			}

			if waterLevel < 0 || waterLevel > 100 {
				t.Errorf("invalid water level: %.2f", waterLevel)
			}

			if code < 1 || code > 3 {
				t.Errorf("invalid status code: %d", code)
			}

			t.Logf("Result → Status: %s | Code: %d | WaterLevel: %.2f",
				status, code, waterLevel)
		})
	}
}