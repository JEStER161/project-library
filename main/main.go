package main

import (
	"project_library/author"
	"project_library/authorization"
	"project_library/book"
	"project_library/config"
	"project_library/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.ConnectDB()      //подключение к базе данных
	defer config.DB.Close() // Закрываем пул соединений при завершении программы

	e := echo.New()
	e.Use(middleware.CORS())

	e.POST("/add_book", book.AddBook, authorization.MiddleWare)             //Добавление новой книги
	e.POST("/add_author", author.AddAuthor, authorization.MiddleWare)       //Добавление нового автора
	e.POST("/add_user", user.AddUser)                                       //Регистрация нового пользователя
	e.POST("/login", user.Sign_in)                                          //Авторизация пользователя
	e.POST("/reserve/:book_id", book.ReserveBook, authorization.MiddleWare) //Бронирование книги
	e.POST("/borrowing_book", book.Borrowing, authorization.MiddleWare)     //Выдача книги

	e.GET("/profile", user.Profile, authorization.MiddleWare)            //Получение профиля пользователя
	e.GET("/check_reserve", user.CheckReserve, authorization.MiddleWare) //Получение всех заявок конкретного пользователя
	e.GET("/all_reserve", user.AllReserve, authorization.MiddleWare)     //Получение ВСЕХ заявок для админа
	e.GET("/all_borrowing", user.AllBorrowing, authorization.MiddleWare) //Получение ВСЕХ выдач книг
	e.GET("/get_book/:book_id", book.GetBook)                            //Получение профиля книги
	e.GET("/search", book.SearchLine)                                    // Поисковая строка
	e.GET("/home", book.AllBook);
	e.GET("/get_author/:author_id", author.GetAuthor)                    //Получение профиля автора

	e.Start(":8080")
}
