package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/edkuzhakhmetov/go_final_project/internal/models"
	"github.com/edkuzhakhmetov/go_final_project/internal/scheduler"
	"github.com/sirupsen/logrus"
)

func (h *Handler) apiPostTaskDone(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Starting api_PostTaskDone")

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
	if task.Repeat != "" {
		now := time.Now()

		nextdate, err := scheduler.NextDate(now, task.Date, task.Repeat)
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
		h.log.WithFields(logrus.Fields{"nextdate": nextdate}).Info("Следующая дата вычислена")

		if task.Date < now.Format("20060102") {
			task.Date = nextdate
		}
		err = h.storage.UpdateDate(ctx, nextdate, id)
		if err != nil {
			h.log.Errorf("Ошибка при обновлении даты задачи: %v", err)
			h.SendResponse(w, models.APIResponse{
				StatusCode: http.StatusInternalServerError,
				Body: models.Body{
					Error: "Ошибка при обновлении даты задачи",
				},
			})
			return
		}
		h.log.WithFields(logrus.Fields{
			"ID":      task.ID,
			"Date":    task.Date,
			"Title":   task.Title,
			"Comment": task.Comment,
			"Repeat":  task.Repeat,
		}).Info("Обновленная задача")

		h.log.Infof("Задача с ID %d обновлена", id)

		h.SendResponse(w, models.APIResponse{
			StatusCode: http.StatusOK,
			Body:       models.Task{},
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
