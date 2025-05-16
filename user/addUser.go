package user

import (
	"project_library/config"
	"project_library/password"
	"project_library/utils"

	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type User struct {
	User_id       int       `json:"user_id"`
	Login         string    `json:"login"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	First_name    string    `json:"first_name"`
	Surname       string    `json:"surname"`
	Patronymic    string    `json:"patronymic"`
	Date_of_birth string    `json:"date_of_birth"`
	Phone         string    `json:"phone"`
	Role          string    `json:"role"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
}

// Регистрация нового пользователя
func AddUser(context echo.Context) error {
	var user User

	if err := context.Bind(&user); err != nil {
		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "Invalid request payload",
		})
	}

	log.Println(user)

	// Проверка на пустой контент
	if user.Login == "" {
		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "Name are required",
		})
	}

	//Хэширование пароля пользователя
	hash_password, err_password := password.HashPassword(user.Password)
	if err_password != nil {
		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "Failed password hashing",
		})
	}

	//Запрос на добавление пользователя в бд
	query := `insert into "library".users(login, email, password_hash, first_name, surname, patronymic, date_of_birth,
										phone, created_at)
									values($1, $2, $3, $4, $5, $6, $7, $8, Now()) returning user_id, created_at, updated_at;`
	err_query := config.DB.QueryRow(context.Request().Context(), query, user.Login, user.Email, hash_password, user.First_name,
		user.Surname, user.Patronymic, user.Date_of_birth, user.Phone).Scan(&user.User_id, &user.Created_at, &user.Updated_at)
	if err_query != nil {

		log.Println(err_query)

		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: "Could not add the user",
		})
	}

	return context.JSON(http.StatusOK, user)
}
