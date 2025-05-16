package user

import (
	"project_library/authorization"
	"project_library/config"
	"project_library/password"
	"project_library/utils"

	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Авторизация
func Sign_in(context echo.Context) error {
	var login authorization.Login

	if err := context.Bind(&login); err != nil {
		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "Invalid request payload",
		})
	}

	if login.Login == "" || login.Password == "" {
		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "Login and password are required",
		})
	}

	query_check := `Select user_id, password_hash, role from "library".users where login = $1;`
	err_login := config.DB.QueryRow(context.Request().Context(), query_check, login.Login).Scan(&login.User_id, &login.Password_from_bd, &login.Role_user)
	if err_login != nil {
		log.Println("Ошибка запроса к БД", err_login)
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: "Could not add the book_author",
		})
	}

	log.Println(login)

	if login.User_id == "" {
		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "incorrect login",
		})
	}

	if err_password := password.CheckPasswordHash(login.Password, login.Password_from_bd); err_password != nil {
		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "incorrect password",
		})
	}

	token, err := authorization.Generate_JWT(login)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: "token generation error",
		})
	}

	return context.JSON(http.StatusOK, authorization.TokenResponse{Token: token})
}
