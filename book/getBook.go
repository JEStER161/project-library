package book

import (
	"log"
	"net/http"
	"project_library/config"
	"project_library/utils"

	"github.com/labstack/echo"
)

func GetBook(context echo.Context) error {
	var book Book //книга, которую будем возвращать
	book_id := context.Param("book_id")

	//Получение данных о книге, кроме ее авторов
	query_getBook := `select title, genre, isbn, total_copies, available_copies, to_char(published_date , 'YYYY-MM-DD'), publisher, description
					from "library".books
					where book_id = $1;`
	rows, err := config.DB.Query(context.Request().Context(), query_getBook, book_id)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: err.Error(),
		})
	}
	defer rows.Close()

	rows.Next()
	if err := rows.Scan(&book.Title, &book.Genre, &book.Isbn, &book.Total_copies, &book.Available_copies, &book.Published_date, &book.Publisher, &book.Description); err != nil {
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: err.Error(),
		})
	}

	//Получение авторов книги
	query_getAuthor := `select a."name"
						from "library".authors as a
						join "library".book_authors as b
						on a.author_id = b.author_id
						where book_id = $1;`

	rows2, err2 := config.DB.Query(context.Request().Context(), query_getAuthor, book_id)
	if err2 != nil {
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: err2.Error(),
		})
	}
	defer rows2.Close()

	var authors []string
	for rows2.Next() {
		var name string
		if err := rows2.Scan(&name); err != nil {
			return context.JSON(http.StatusInternalServerError, utils.Response{
				Status:  "Error",
				Message: err.Error(),
			})
		}
		authors = append(authors, name)
	}

	log.Println(authors)
	author_result := ""

	//Соедиение всех полученных авторов в одну строку
	for _, elem := range authors {
		author_result = author_result + ", " + elem
	}

	book.Author = author_result

	return context.JSON(http.StatusOK, book)
}
