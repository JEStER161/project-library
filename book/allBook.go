package book

import (
	"project_library/config"
	"project_library/utils"

	"net/http"

	"github.com/labstack/echo/v4"
)

func AllBook(context echo.Context) error {
	var books []Book
	var book Book

	query := `select book_id, title, total_copies, available_copies, cover_image from "library".books;`
	rows_books, err_books := config.DB.Query(context.Request().Context(), query)
	defer rows_books.Close()

	if err_books != nil {
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: err_books.Error(),
		})
	}

	for rows_books.Next() {
		err := rows_books.Scan(&book.Book_id, &book.Title, &book.Total_copies, &book.Available_copies, &book.Cover_image)
		if err != nil {
			return context.JSON(http.StatusInternalServerError, utils.Response{
				Status:  "Error",
				Message: err.Error(),
			})
		}
		books = append(books, book)
	}

	return context.JSON(http.StatusOK, books)
}
