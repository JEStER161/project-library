package user

import (
	"net/http"
	"project_library/config"
	"project_library/utils"

	"github.com/labstack/echo"
)

func Profile(context echo.Context) error {
	var user User

	user_id := context.Get("user_id").(string)

	query_profile := `select email, first_name, surname, patronymic, to_char(date_of_birth, 'YYYY-MM-DD'), phone from "library".users where user_id = $1;`
	rows, err_profile := config.DB.Query(context.Request().Context(), query_profile, user_id)
	if err_profile != nil {
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: err_profile.Error(),
		})
	}
	defer rows.Close()

	rows.Next()
	if err := rows.Scan(&user.Email, &user.First_name, &user.Surname, &user.Patronymic, &user.Date_of_birth, &user.Phone); err != nil {
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: err.Error(),
		})
	}

	return context.JSON(http.StatusOK, user)
}
