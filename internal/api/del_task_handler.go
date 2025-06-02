package api

import (
	"database/sql"
	"net/http"

	"github.com/edkuzhakhmetov/go_final_project/internal/models"
	"github.com/edkuzhakhmetov/go_final_project/internal/scheduler"
)

func (h *Handler) delTaskHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Starting delTaskHandler action")

	ctx := r.Context()
	query := r.URL.Query()

	reqId := query.Get("id")

	id, err := scheduler.ValidateTaskID(reqId)
	if err != nil {
		h.log.Warnf("Ошибка валидации ID задачи: %v", err)
		h.SendResponse(w, models.APIResponse{
			StatusCode: http.StatusBadRequest,
			Body: models.Body{
				Error: err.Error(),
			},
		})
		return
	}

	_, err = h.storage.GetTask(ctx, id)
	if err != nil {

		if err == sql.ErrNoRows {
			h.log.Infof("Задача с ID %d не найдена", id)
			h.SendResponse(w, models.APIResponse{
				StatusCode: http.StatusNotFound,
				Body: models.Body{
					Error: "Задача не найдена",
				},
			})

			return
		}

		h.log.Errorf("Ошибка при получении задачи: %v", err)
		h.SendResponse(w, models.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Body: models.Body{
				Error: "Ошибка при получении задачи",
			},
		})

		return
	}
	err = h.storage.DeleteTask(ctx, id)
	if err != nil {
		h.log.Errorf("Ошибка при обновлении задачи: %v", err)
		h.SendResponse(w, models.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Body: models.Body{
				Error: "Ошибка при обновлении задачи",
			},
		})

		return
	}

	h.log.Infof("Задача с ID %d обновлена", id)
	h.SendResponse(w, models.APIResponse{
		StatusCode: http.StatusOK,
		Body:       models.Task{},
	})

}
