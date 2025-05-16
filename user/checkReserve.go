package user

import (
	"project_library/config"
	"project_library/utils"

	"net/http"

	"github.com/labstack/echo/v4"
)

type reserve struct {
	Book_title         string `json:"title"`
	Isbn               string `json:"isbn"`
	Reservation_date   string `json:"res_date"`
	End_of_reservation string `json:"end_res"`
	Status             string `json:"status"`
}

func CheckReserve(context echo.Context) error {
	user_id := context.Get("user_id").(string)
	var reserves []reserve
	var reserve reserve

	query_check_reserve := `select b.title, b.isbn, to_char(reservation_date, 'YYYY-MM-DD'), to_char(end_of_reserve, 'YYYY-MM-DD'), status
							from  "library".reservations as r
							join "library".books as b
							on r.book_id = b.book_id
							where user_id = $1;`

	rows, err := config.DB.Query(context.Request().Context(), query_check_reserve, user_id)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: err.Error(),
		})
	}

	for rows.Next() {
		if err := rows.Scan(&reserve.Book_title, &reserve.Isbn, &reserve.Reservation_date, &reserve.End_of_reservation, &reserve.Status); err != nil {
			return context.JSON(http.StatusInternalServerError, utils.Response{
				Status:  "Error",
				Message: err.Error(),
			})
		}

		reserves = append(reserves, reserve)
	}

	return context.JSON(http.StatusOK, reserves)
}
