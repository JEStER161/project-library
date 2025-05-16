package user

import (
	"log"
	"project_library/config"
	"project_library/utils"

	"net/http"

	"github.com/labstack/echo/v4"
)

type allReserve struct {
	Reserve_id         string `json:"reserve_id"`
	Book_id            string `json:"book_id"`
	Book_title         string `json:"title"`
	Isbn               string `json:"isbn"`
	User_id            string `json:"user_id"`
	Name               string `json:"name"`
	Reservation_date   string `json:"res_date"`
	End_of_reservation string `json:"end_res"`
	Status             string `json:"status"`
}

func AllReserve(context echo.Context) error {
	var reserves []allReserve
	var reserve allReserve

	if context.Get("role_user").(string) != "admin" {

		log.Println(context.Get("role").(string))

		return context.JSON(http.StatusBadRequest, utils.Response{
			Status:  "Error",
			Message: "You don't have sufficient rights",
		})
	}

	query_all_reserve := `select r.reserve_id,  b.book_id, b.title, b.isbn, u.user_id, u.surname || ' ' || u.first_name || ' ' || u.patronymic,
							to_char(r.reservation_date, 'YYYY-MM-DD'), to_char(r.end_of_reserve, 'YYYY-MM-DD'), r.status
							from "library".reservations as r
							join "library".books as b
							on r.book_id = b.book_id
							join "library".users as u
							on r.user_id = u.user_id;`

	rows, err := config.DB.Query(context.Request().Context(), query_all_reserve)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "Error",
			Message: err.Error(),
		})
	}

	for rows.Next() {
		err := rows.Scan(&reserve.Reserve_id, &reserve.Book_id, &reserve.Book_title, &reserve.Isbn, &reserve.User_id, &reserve.Name,
			&reserve.Reservation_date, &reserve.End_of_reservation, &reserve.Status)
		if err != nil {
			return context.JSON(http.StatusInternalServerError, utils.Response{
				Status:  "Error",
				Message: err.Error(),
			})
		}

		reserves = append(reserves, reserve)
	}

	return context.JSON(http.StatusOK, reserves)
}
