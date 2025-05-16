package user

// import (
// 	"project_library/config"

// 	"github.com/labstack/echo"
// )

// func CheckRole(context echo.Context, user_id string) string {
// 	var role string
// 	query_check := `select role from "library".users where user_id = $1;`
// 	err := config.DB.QueryRow(context.Request().Context(), query_check, user_id).Scan(&role)

// 	if err != nil {
// 		return err.Error()
// 	}

// 	return role
// }
