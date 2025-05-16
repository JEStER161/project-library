package book

import (
	"net/http"
	"project_library/config"
	"project_library/utils"

	"github.com/labstack/echo/v4"
)

type Borrowings struct {
	Borrowing_id int    `json:"borrowing_id"`
	User_id      int    `json:"user_id"`
	Book_id      int    `json:"book_id"`
	Borrow_date  string `json:"borrow_date"`
	Due_date     string `json:"due_date"`
	Return_date  string `json:"return_date"`
	Status       string `json:"status"`
}

func Borrowing(context echo.Context) error {
	var borrowing Borrowings

	if context.Get("role_user") != "admin" {
		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "You don't have sufficient rights",
		})
	}

	if err := context.Bind(&borrowing); err != nil {
		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "Invalid request payload",
		})
	}

	query := `insert into "library".borrowings(user_id, book_id, borrow_date, status)
				values ($1, $2, NOW(), $3) returning to_char(due_date, 'YYYY-MM-DD'), borrowing_id, to_char(borrow_date, 'YYYY-MM-DD');`

	err := config.DB.QueryRow(context.Request().Context(), query, borrowing.User_id, borrowing.Book_id, borrowing.Status).Scan(&borrowing.Due_date, &borrowing.Borrowing_id,
		&borrowing.Borrow_date)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: err.Error(),
		})
	}

	return context.JSON(http.StatusOK, borrowing)
}
