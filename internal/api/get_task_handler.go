package api

import (
	"database/sql"
	"net/http"

	"github.com/edkuzhakhmetov/go_final_project/internal/models"
	"github.com/edkuzhakhmetov/go_final_project/internal/scheduler"
)

func (h *Handler) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Starting getTaskHandler action")

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

	task, err := h.storage.GetTask(ctx, id)
	if err != nil {

		if err == sql.ErrNoRows || task.ID == "" {
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

	h.log.Infof("Задача с ID %d получена", id)
	h.SendResponse(w, models.APIResponse{
		StatusCode: http.StatusCreated,
		Body: models.Task{
			ID:      task.ID,
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		},
	})
}
