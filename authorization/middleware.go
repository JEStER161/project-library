package authorization

import (
	"errors"
	"log"
	"net/http"
	"project_library/utils"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	User_id   string `json:"user_id"`
	Role_user string `json:"role_user"`
	jwt.RegisteredClaims
}

func MiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		//Получаем токен из зголовка Authorization
		authHeader := context.Request().Header.Get("Authorization")

		//Проверка на пустоту
		if authHeader == "" {
			return context.JSON(http.StatusUnauthorized, utils.Response{
				Status:  "Error",
				Message: "Missing token",
			})
		}

		// Токен передается в формате "Bearer <token>", извлекаем сам токен
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return context.JSON(http.StatusUnauthorized, utils.Response{
				Status:  "Error",
				Message: "Invalid token format",
			})
		}
		tokenString := parts[1]

		//Разбираем токен
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Проверяем метод подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("unexpected signing method: %v", token.Header["alg"])
				return nil, errors.New("invalid token signing method")
			}
			return []byte(secretKey), nil
		})

		//// Проверяем ошибки
		if err != nil || !token.Valid {
			return context.JSON(http.StatusUnauthorized, utils.Response{
				Status:  "Error",
				Message: "Invalid or expired token",
			})
		}

		// Проверяем срок действия токена
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return context.JSON(http.StatusUnauthorized, utils.Response{
				Status:  "Error",
				Message: "Token expired",
			})
		}

		// Передаем user_id в контекст запроса
		context.Set("user_id", claims.User_id)
		context.Set("role_user", claims.Role_user)

		// Вызываем следующий обработчик
		return next(context)
	}
}
