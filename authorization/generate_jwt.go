package authorization

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("my_secret_key")

type TokenResponse struct {
	Token string `json:"token"`
}

type Login struct {
	Login            string `json:"login"`
	Password         string `json:"password"`
	Role_user        string `json:"role_user"`
	Password_from_bd string `json:"password_from_bd"`
	User_id          string `json:"user_id"`
}

// генерация jwt
func Generate_JWT(login Login) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   login.User_id,
		"role_user": login.Role_user,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Токен истекает через 24 часа
	})
	// Подписываем токен секретным ключом
	return token.SignedString(secretKey)
}
