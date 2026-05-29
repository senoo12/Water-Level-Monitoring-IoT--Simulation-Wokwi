package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY = []byte("super-secret-key")

func GenerateToken(deviceID string) (string, error)  {
	claims := jwt.MapClaims{
		"device_id": deviceID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SECRET_KEY)
}