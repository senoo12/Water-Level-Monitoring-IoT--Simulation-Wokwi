package repository

type SensorRepository interface {
	SendData(water, temp, pressure, distance float64, status int) error
}