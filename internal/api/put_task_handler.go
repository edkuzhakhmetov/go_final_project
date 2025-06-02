package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/edkuzhakhmetov/go_final_project/internal/models"
	"github.com/edkuzhakhmetov/go_final_project/internal/scheduler"
	"github.com/sirupsen/logrus"
)

func (h *Handler) apiPutTask(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Starting api_PutTask action")

	ctx := r.Context()

	task, err := h.getTaskFromBody(r)

	if err != nil {
		h.SendResponse(w, models.APIResponse{
			StatusCode: http.StatusBadRequest,
			Body: models.Body{
				Error: "Не удалось прочитать тело запроса",
			},
		})

		return
	}

	id, err := scheduler.ValidateTaskID(task.ID)
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

	err = scheduler.ValidateTask(task)
	if err != nil {
		h.log.Error(err)
		h.SendResponse(w, models.APIResponse{
			StatusCode: http.StatusBadRequest,
			Body: models.Body{
				Error: err.Error(),
			},
		})

		return
	}

	if task.Date == "" {
		h.log.Error("Поле Date обязательно")
		h.SendResponse(w, models.APIResponse{
			StatusCode: http.StatusBadRequest,
			Body: models.Body{
				Error: "Поле Date обязательно",
			},
		})

		return
	}
	if _, err = time.Parse(scheduler.Format, task.Date); err != nil {

		h.log.Errorf("Ошибка при парсинге даты: %v", err)
		h.SendResponse(w, models.APIResponse{
			StatusCode: http.StatusBadRequest,
			Body: models.Body{
				Error: "Некорректный формат даты",
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

	if task.Repeat != "" {
		_, err := scheduler.NextDate(time.Now(), task.Date, task.Repeat)

		if err != nil {
			h.log.WithFields(logrus.Fields{
				"ID":      task.ID,
				"Date":    task.Date,
				"Title":   task.Title,
				"Comment": task.Comment,
				"Repeat":  task.Repeat,
			}).Errorf("Ошибка при вычислении следующей даты: %v", err)

			h.SendResponse(w, models.APIResponse{
				StatusCode: http.StatusBadRequest,
				Body: models.Body{
					Error: "Ошибка при вычислении следующей даты. Проверьте формат даты или правило повторения",
				},
			})

			return
		}

	}
	err = h.storage.UpdateTask(ctx, task)
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

	h.log.Infof("Задача с ID %d обновлена", task.ID)
	h.SendResponse(w, models.APIResponse{
		StatusCode: http.StatusOK,
		Body:       models.Task{},
	})
}
