package book

import (
	"project_library/author"
	"project_library/config"
	"project_library/utils"

	"net/http"

	"github.com/labstack/echo/v4"
)

func SearchLine(context echo.Context) error {
	searchline := context.QueryParam("q")
	var books []Book
	var book Book
	var authors []author.Author
	var author author.Author

	query_search_author := `select name from "library".authors where name ilike('%' || $1 || '%');`
	rows_authors, err_authors := config.DB.Query(context.Request().Context(), query_search_author, searchline)

	query_search_book := `select title, total_copies, available_copies, cover_image from "library".books
							where title ilike ('%' || $1 || '%');`
	rows_books, err_books := config.DB.Query(context.Request().Context(), query_search_book, searchline)
	defer rows_authors.Close()
	defer rows_books.Close()

	if err_authors != nil {
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: err_authors.Error(),
		})
	}

	if err_books != nil {
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: err_books.Error(),
		})
	}

	for rows_books.Next() {
		err := rows_books.Scan(&book.Title, &book.Total_copies, &book.Available_copies, &book.Cover_image)
		if err != nil {
			return context.JSON(http.StatusInternalServerError, utils.Response{
				Status:  "Error",
				Message: err.Error(),
			})
		}
		books = append(books, book)
	}

	for rows_authors.Next() {
		err := rows_authors.Scan(&author.Name)
		if err != nil {
			return context.JSON(http.StatusInternalServerError, utils.Response{
				Status:  "Error",
				Message: err.Error(),
			})
		}
		authors = append(authors, author)
	}

	return context.JSON(http.StatusOK, map[string]interface{}{
		"books":   books,
		"authors": authors,
	})
}
