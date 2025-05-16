package book

import (
	"log"
	"net/http"
	"project_library/config"
	"project_library/utils"

	"github.com/labstack/echo/v4"
)

func ReserveBook(context echo.Context) error {
	user_id := context.Get("user_id").(string) //получили id пользователя из context
	book_id := context.Param("book_id")     //получили id книги из адресса

	log.Println(user_id, book_id)

	//Запрос на добавление брони книги
	query_reserve := `insert into "library".reservations(user_id, book_id) values($1, $2)`

	_, err_reserve := config.DB.Query(context.Request().Context(), query_reserve, user_id, book_id)
	if err_reserve != nil {
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: err_reserve.Error(),
		})
	}

	return context.JSON(http.StatusOK, utils.Response{
		Status:  "Ok",
		Message: "Бронь одобрена",
	})
}
