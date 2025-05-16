package user

import (
	"log"
	"net/http"
	"project_library/config"
	"project_library/utils"

	"github.com/labstack/echo/v4"
)

type allBorrowings struct {
	Borrowing_id int    `json:"borrowing_id"`
	Book_id      int    `json:"book_id"`
	Book_title   string `json:"title"`
	User_id      int    `json:"user_id"`
	Name         string `json:"name"`
	Borrow_date  string `json:"borrow_date"`
	Due_date     string `json:"due_date"`
	Return_date  string `json:"return_date"`
	Status       string `json:"status"`
}

func AllBorrowing(context echo.Context) error {
	var borrowings []allBorrowings
	var borrowing allBorrowings

	if context.Get("role_user").(string) != "admin" {

		log.Println(context.Get("role").(string))

		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "You don't have sufficient rights",
		})
	}

	query := `select b.borrowing_id, b.book_id, bo.title, b.user_id, u.surname || ' ' || u.first_name || ' ' || u.patronymic,
				to_char(b.borrow_date, 'YYYY-MM-DD'), to_char(b.due_date, 'YYYY-MM-DD'), to_char(b.return_date , 'YYYY-MM-DD'),
				b.status
				from "library".borrowings as b
				join "library".books as bo
				on bo.book_id = b.book_id
				join "library".users as u
				on u.user_id = b.user_id;`

	rows, err := config.DB.Query(context.Request().Context(), query)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: err.Error(),
		})
	}

	for rows.Next() {
		err := rows.Scan(&borrowing.Borrowing_id, &borrowing.Book_id, &borrowing.Book_title, &borrowing.User_id, &borrowing.Name, &borrowing.Borrow_date, &borrowing.Due_date,
			&borrowing.Return_date, &borrowing.Status)
		if err != nil {
			return context.JSON(http.StatusInternalServerError, utils.Response{
				Status:  "Error",
				Message: err.Error(),
			})
		}

		borrowings = append(borrowings, borrowing)
	}

	return context.JSON(http.StatusOK, borrowings)
}
