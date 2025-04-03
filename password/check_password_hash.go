package password

import (
	"golang.org/x/crypto/bcrypt"
)

// Проверка пароля
func CheckPasswordHash(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
