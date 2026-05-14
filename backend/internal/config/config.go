package config

import "os"

type Config struct {
	ThingSpeakAPIKey string
	GsheetWebHook string
}

func LoadConfig() *Config {
	return &Config{
		ThingSpeakAPIKey: os.Getenv("THINGSPEAK_API_KEY"),
		GsheetWebHook: os.Getenv("GSHEET_WEBHOOK"),
	}
}