package main

import (
	"project_library/author"
	"project_library/authorization"
	"project_library/book"
	"project_library/config"
	"project_library/user"

	"github.com/labstack/echo"
)

func main() {
	config.ConnectDB()      //подключение к базе данных
	defer config.DB.Close() // Закрываем пул соединений при завершении программы

	e := echo.New()

	e.POST("/add_book", book.AddBook) //Добавление новой книги
	e.POST("/add_author", author.AddAuthor) //Добавление нового автора
	e.POST("/add_user", user.AddUser) //Регистрация нового пользователя
	e.POST("/login", user.Sign_in) //Авторизация пользователя
	e.POST("/reserve/:book_id", book.ReserveBook, authorization.MiddleWare) //Бронирование книги

	e.GET("/profile", user.Profile, authorization.MiddleWare) //Получение профиля пользователя
	e.GET("/get_book/:book_id", book.GetBook) //Получение профиля книги

	e.Start(":8080")
}
