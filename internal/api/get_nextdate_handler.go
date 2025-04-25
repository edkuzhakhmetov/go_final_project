package api

import (
	"net/http"
	"time"

	"github.com/edkuzhakhmetov/go_final_project/internal/models"
	"github.com/edkuzhakhmetov/go_final_project/internal/scheduler"
)

func (h *Handler) NextDayHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Starting NextDayHandler action")

	query := r.URL.Query()
	nowStr := query.Get("now")
	dateStr := query.Get("date")
	repeat := query.Get("repeat")

	var err error

	if nowStr == "" || dateStr == "" || repeat == "" {
		h.log.Error("Параметры now, date и repeat обязательны")

		h.SendResponse(w, models.APIResponse{
			StatusCode: http.StatusBadRequest,
			Body: models.Body{
				Error: "Параметры now, date и repeat обязательны",
			},
		})

		return
	}

	var now time.Time

	if now, err = time.Parse(scheduler.Format, nowStr); err != nil {
		h.log.Errorf("Ошибка при парсинге даты now: %v", err)

		h.SendResponse(w, models.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Body: models.Body{
				Error: "An unexpected error occurred",
			},
		})

		return
	}
	date, err := scheduler.NextDate(now, dateStr, repeat)
	if err != nil {
		h.log.Error("Параметры now, date и repeat обязательны")

		h.SendResponse(w, models.APIResponse{
			StatusCode: http.StatusBadRequest,
			Body: models.Body{
				Message: "Параметры now, date или repeat заполнены некорректно",
			},
		})

		return
	}
	w.Write([]byte(date))

}
