package api

import (
	"database/sql"
	"net/http"

	"github.com/edkuzhakhmetov/go_final_project/internal/models"
)

func (h *Handler) getTasksHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Starting getTasksHandler action")

	ctx := r.Context()
	limit := 50

	resp := make([]*models.Task, 0)

	tasks, err := h.storage.GetTasks(ctx, limit)
	if err != nil {

		if err == sql.ErrNoRows {
			h.log.Infof("Задачи не найдены")
			h.SendResponse(w, models.APIResponse{
				StatusCode: http.StatusNotFound,
				Body: models.Body{
					Error: "Задачи не найдены",
				},
			})

			return
		}

		h.log.Errorf("Ошибка при получении списка задач: %v", err)
		h.SendResponse(w, models.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Body: models.Body{
				Error: "Ошибка при получении списка задач",
			},
		})

		return
	}

	if len(tasks) > 0 {
		resp = tasks[:]
	}

	h.SendResponse(w, models.APIResponse{
		StatusCode: http.StatusCreated,
		Body: models.TasksResp{
			Tasks: resp,
		},
	})
}
