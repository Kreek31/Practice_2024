package handlers

import (
	"myproject/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Секрет для подписи JWT
var jwtSecret = []byte("my_secret_key")

// Генерация JWT для пользователя
func GenerateJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID.Hex(),
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	return token.SignedString(jwtSecret)
}
