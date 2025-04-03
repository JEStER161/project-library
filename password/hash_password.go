package password

import (
	"golang.org/x/crypto/bcrypt"
)

// Хэширование пароля
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
