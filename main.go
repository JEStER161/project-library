package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
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

type User struct {
	User_id       int       `json:"user_id"`
	Login         string    `json:"login"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	First_name    string    `json:"first_name"`
	Surname       string    `json:"surname"`
	Patronymic    string    `json:"patronymic"`
	Date_of_birth string    `json:"date_of_birth"`
	Phone         string    `json:"phone"`
	Role          string    `json:"role"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
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
	log.Println("Подключение к БД установлено!")
}

func separation_authors(authors string) []string {
	res := strings.Split(authors, ",")
	return res
}

// Хэширование пароля
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// Проверка пароля
func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func AddUser(context echo.Context) error {
	var user User

	if err := context.Bind(&user); err != nil {
		return context.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid request payload",
		})
	}

	log.Println(user)

	// Проверка на пустой контент
	if user.Login == "" {
		return context.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Name are required",
		})
	}

	//Хэширование пароля пользователя
	hash_password, err_password := HashPassword(user.Password)
	if err_password != nil {
		return context.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Failed password hashing",
		})
	}

	//Запрос на добавление пользователя в бд
	query := `insert into "library".users(login, email, password_hash, first_name, surname, patronymic, date_of_birth,
										phone, "role", created_at)
									values($1, $2, $3, $4, $5, $6, $7, $8, $9, Now()) returning user_id, created_at, updated_at;`
	err_query := DB.QueryRow(context.Request().Context(), query, user.Login, user.Email, hash_password, user.First_name,
		user.Surname, user.Patronymic, user.Date_of_birth, user.Phone, user.Role).Scan(&user.User_id, &user.Created_at, &user.Updated_at)
	if err_query != nil {

		log.Println(err_query)

		return context.JSON(http.StatusInternalServerError, Response{
			Status:  "Error",
			Message: "Could not add the user",
		})
	}

	return context.JSON(http.StatusOK, user)
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

	authors := separation_authors(book.Author)

	log.Println(book)

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
		log.Println(err_book)
		tx.Rollback(context.Request().Context())
		return context.JSON(http.StatusInternalServerError, Response{
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
			return context.JSON(http.StatusInternalServerError, Response{
				Status:  "Error",
				Message: "There is no author with name: " + authors[i],
			})
		}

		query_book_author := `insert into "library".book_authors (book_id, author_id) values($1, $2);`
		_, err_book_author := tx.Exec(context.Request().Context(), query_book_author, book.Book_id, author_id)
		if err_book_author != nil {
			tx.Rollback(context.Request().Context())
			log.Println("Ошибка запроса к БД", err_book_author)
			return context.JSON(http.StatusInternalServerError, Response{
				Status:  "Error",
				Message: "Could not add the book_author",
			})
		}
		log.Println("Запись добавлена")
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

func AddAuthor(context echo.Context) error {
	var author Author

	if err := context.Bind(&author); err != nil {
		return context.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid request payload",
		})
	}

	log.Println(author)

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
		log.Println("Ошибка запроса к БД", err_author)
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
	e.POST("/add_user", AddUser)

	e.Start(":8080")

}
