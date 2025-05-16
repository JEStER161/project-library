package author

import (
	"project_library/config"
	"project_library/utils"

	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAuthor(context echo.Context) error {
	author_id := context.Param("author_id")
	var author Author

	query_get_author := `select name,  to_char(date_of_birth , 'YYYY-MM-DD'), country, bio from "library".authors where author_id = $1;`

	row, err_author := config.DB.Query(context.Request().Context(), query_get_author, author_id)
	if err_author != nil {
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: err_author.Error(),
		})
	}
	defer row.Close()

	row.Next()
	if err := row.Scan(&author.Name, &author.Date_of_birth, &author.Country, &author.Bio); err != nil {
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: err.Error(),
		})
	}
	log.Println(author)

	return context.JSON(http.StatusOK, author)
}
