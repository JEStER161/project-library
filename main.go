package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo"

	"strings"
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

type Author struct {
	Author_id     int       `json:"author_id"`
	Name          string    `json:"name"`
	Date_of_birth string    `json:"date_of_birth"`
	Country       string    `json:"country"`
	Bio           string    `json:"bio"`
	Created_at    time.Time `json:"created_at"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var DB *pgxpool.Pool

func ConnectDB() {
	dsn := "host=localhost user=postgres password=password dbname=postgres port=1234 sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	DB = pool
	fmt.Println("Подключение к БД установлено!")
}

func separation_authors(authors string) []string {
	res := strings.Split(authors, ",")
	return res
}

func AddBook(context echo.Context) error {
	var book Book
	tx, eror := DB.Begin(context.Request().Context())
	if eror != nil {
		return context.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Transaction failed",
		})
	}

	// Привязка JSON к структуре
	if err := context.Bind(&book); err != nil {
		return context.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid request payload",
		})
	}

	fmt.Println(book.Author)
	authors := separation_authors(book.Author)
	fmt.Println(authors)

	fmt.Println(book)

	// Проверка на пустой контент
	if book.Title == "" {
		return context.JSON(http.StatusBadRequest, Response{
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
		fmt.Println(err_book)
		tx.Rollback(context.Request().Context())
		return context.JSON(http.StatusInternalServerError, Response{
			Status:  "Error",
			Message: "Could not add the book",
		})
	}

	// Запрос на добавление связи между авторами и книгой
	for i := 0; i < len(authors); i++ {
		fmt.Println(authors[i])
		var author_id int
		err_1 := tx.QueryRow(context.Request().Context(), `select author_id from "library".authors where name = $1;`, authors[i]).Scan(&author_id)
		if err_1 != nil {
			fmt.Println(err_1)
			tx.Rollback(context.Request().Context())
			return context.JSON(http.StatusInternalServerError, Response{
				Status:  "Error",
				Message: "There is no author with name: " + authors[i],
			})
		}
		fmt.Println(book.Book_id, author_id)

		query_book_author := `insert into "library".book_authors (book_id, author_id) values($1, $2);`
		_, err_book_author := tx.Exec(context.Request().Context(), query_book_author, book.Book_id, author_id)
		if err_book_author != nil {
			tx.Rollback(context.Request().Context())
			fmt.Println(err_book_author)
			return context.JSON(http.StatusInternalServerError, Response{
				Status:  "Error",
				Message: "Could not add the book_author",
			})
		}
		fmt.Println("Запись добавлена")
	}

	error_3 := tx.Commit(context.Request().Context())
	if error_3 != nil {
		return context.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Commit failed",
		})
	}
	// Успешный ответ
	return context.JSON(http.StatusOK, book)
}

/*func Find_Book_from_author(context echo.Context) error {
	var slice []Book
	author_name := context.Param("name")

	query := `Select `



}*/

func AddAuthor(context echo.Context) error {
	var author Author

	if err := context.Bind(&author); err != nil {
		return context.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid request payload",
		})
	}

	fmt.Println(author)

	// Проверка на пустой контент
	if author.Name == "" {
		return context.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Name are required",
		})
	}

	//Запрос на добавление автора в бд
	query_author := `insert into "library".authors (name, date_of_birth, country, bio, created_at)
					values($1, $2, $3, $4, NOW()) returning author_id;`
	err_author := DB.QueryRow(context.Request().Context(), query_author, author.Name,
		author.Date_of_birth, author.Country, author.Bio).Scan(&author.Author_id)

	if err_author != nil {
		fmt.Println(err_author)
		return context.JSON(http.StatusInternalServerError, Response{
			Status:  "Error",
			Message: "Could not add the author",
		})
	}

	return context.JSON(http.StatusOK, author)
}

func main() {
	ConnectDB()      //подключение к базе данных
	defer DB.Close() // Закрываем пул соединений при завершении программы

	e := echo.New()

	e.POST("/add_book", AddBook)
	e.POST("/add_author", AddAuthor)
	//e.GET("/find_Book_from_author/:name", Find_Book_from_author)

	e.Start(":8080")

}
