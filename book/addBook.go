package book

import (
	"project_library/config"
	"project_library/utils"

	"github.com/labstack/echo"

	"log"
	"net/http"
	"strings"
	"time"
)

type Book struct {
	Book_id          int       `json:"book_id"`
	Title            string    `json:"title"`
	Author           string    `json:"author"`
	Genre            string    `json:"genre"`
	Isbn             string    `json:"isbn"`
	Total_copies     int       `json:"total_copies"`
	Available_copies int       `json:"available_copies"`
	Published_date   string    `json:"published_date"`
	Publisher        string    `json:"publisher"`
	Description      string    `json:"description"`
	Cover_image      string    `json:"cover_image"`
	Created_at       time.Time `json:"created_at"`
	Updated_at       time.Time `json:"updated_at"`
}

// Разделение строки с авторами
func separation_authors(authors string) []string {
	res := strings.Split(authors, ",")
	return res
}

// Добавление новой книги
func AddBook(context echo.Context) error {
	var book Book
	tx, eror := config.DB.Begin(context.Request().Context())
	if eror != nil {
		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "Transaction failed",
		})
	}

	// Привязка JSON к структуре
	if err := context.Bind(&book); err != nil {
		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "Invalid request payload",
		})
	}

	authors := separation_authors(book.Author)

	log.Println(book)

	// Проверка на пустой контент
	if book.Title == "" {
		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "Title are required",
		})
	}
	//Запрос на добавление книги
	query_book := `INSERT INTO "library".books (title, genre, isbn, total_copies, available_copies,
				published_date, publisher, description, cover_image, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW()) RETURNING book_id, created_at, updated_at;`

	err_book := tx.QueryRow(context.Request().Context(), query_book, book.Title, book.Genre, book.Isbn, book.Total_copies,
		book.Available_copies, book.Published_date, book.Publisher, book.Description,
		book.Cover_image).Scan(&book.Book_id, &book.Created_at, &book.Updated_at)

	if err_book != nil {
		log.Println(err_book)
		tx.Rollback(context.Request().Context())
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: "Could not add the book",
		})
	}

	// Запрос на добавление связи между авторами и книгой
	for i := 0; i < len(authors); i++ {
		var author_id int
		err_1 := tx.QueryRow(context.Request().Context(), `select author_id from "library".authors where name = $1;`, authors[i]).Scan(&author_id)
		if err_1 != nil {
			log.Println(err_1)
			tx.Rollback(context.Request().Context())
			return context.JSON(http.StatusInternalServerError, utils.Response{
				Status:  "Error",
				Message: "There is no author with name: " + authors[i],
			})
		}

		query_book_author := `insert into "library".book_authors (book_id, author_id) values($1, $2);`
		_, err_book_author := tx.Exec(context.Request().Context(), query_book_author, book.Book_id, author_id)
		if err_book_author != nil {
			tx.Rollback(context.Request().Context())
			log.Println("Ошибка запроса к БД", err_book_author)
			return context.JSON(http.StatusInternalServerError, utils.Response{
				Status:  "Error",
				Message: "Could not add the book_author",
			})
		}
		log.Println("Запись добавлена")
	}

	error_3 := tx.Commit(context.Request().Context())
	if error_3 != nil {
		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "Commit failed",
		})
	}
	// Успешный ответ
	return context.JSON(http.StatusOK, book)
}
