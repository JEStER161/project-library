package author

import (
	"project_library/config"
	"project_library/utils"

	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Author struct {
	Author_id     int       `json:"author_id"`
	Name          string    `json:"name"`
	Date_of_birth string    `json:"date_of_birth"`
	Country       string    `json:"country"`
	Bio           string    `json:"bio"`
	Created_at    time.Time `json:"created_at"`
}

// Добавление нового автора
func AddAuthor(context echo.Context) error {
	var author Author

	if context.Get("role_user").(string) != "admin"{
		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "You don't have sufficient rights",
		})
	}

	if err := context.Bind(&author); err != nil {
		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "Invalid request payload",
		})
	}

	log.Println(author)

	// Проверка на пустой контент
	if author.Name == "" {
		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "Name are required",
		})
	}

	//Запрос на добавление автора в бд
	query_author := `insert into "library".authors (name, date_of_birth, country, bio, created_at)
					values($1, $2, $3, $4, NOW()) returning author_id;`
	err_author := config.DB.QueryRow(context.Request().Context(), query_author, author.Name,
		author.Date_of_birth, author.Country, author.Bio).Scan(&author.Author_id)

	if err_author != nil {
		log.Println("Ошибка запроса к БД", err_author)
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: "Could not add the author",
		})
	}

	return context.JSON(http.StatusOK, author)
}
