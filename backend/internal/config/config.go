package config

import "os"

type Config struct {
	ThingSpeakAPIKey string
}

func LoadConfig() *Config {
	return &Config{
		ThingSpeakAPIKey: os.Getenv("THINGSPEAK_API_KEY"),
	}
}