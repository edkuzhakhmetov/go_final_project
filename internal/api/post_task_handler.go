package api

import (
	"net/http"
	"time"

	"github.com/edkuzhakhmetov/go_final_project/internal/models"
	"github.com/edkuzhakhmetov/go_final_project/internal/scheduler"
	"github.com/sirupsen/logrus"
)

func (h *Handler) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Starting PostTaskHandler action")

	ctx := r.Context()

	reqTask, err := h.getTaskFromBody(r)

	if err != nil {
		h.SendResponse(w, models.APIResponse{
			StatusCode: http.StatusBadRequest,
			Body: models.Body{
				Error: "Не удалось прочитать тело запроса",
			},
		})

		return
	}

	//Нужно десериализовать полученный в запросе JSON в переменную var task db.Task.

	//не указан заголовок задачи;
	//Проверить, что поле task.Title не пустое.
	err = scheduler.ValidateTask(reqTask)
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

	//не указана дата задачи;
	//если дата не указана, то используем текущую дату
	now := time.Now()

	if reqTask.Date == "" {
		//если task.Date пустая строка, то присваиваем ему текущее время now.Format("20060102");
		reqTask.Date = now.Format(scheduler.Format)
		h.log.Infof("Date is nil, using current date: %s", reqTask.Date)
	}

	var date time.Time

	//дата представлена в формате, отличном от 20060102;
	if date, err = time.Parse(scheduler.Format, reqTask.Date); err != nil {
		//проверяем, что в task.Date указана корректная дата t, err := time.Parse("20060102", task.Date). t нам ещё пригодится;
		h.log.Errorf("Ошибка при парсинге даты: %v", err)
		h.SendResponse(w, models.APIResponse{
			StatusCode: http.StatusBadRequest,
			Body: models.Body{
				Error: "Ошибка при парсинге даты",
			},
		})

		return
	}

	var nextdate string

	//при указанном правиле повторения вам нужно вычислить
	//если определён task.Repeat, то проверяем корректность правила и заодно получаем следующую дату next, err = NextDate(now, task.Date, task.Repeat);

	if reqTask.Repeat != "" {
		nextdate, err = scheduler.NextDate(now, reqTask.Date, reqTask.Repeat)

		//правило повторения указано в неправильном формате.
		if err != nil {
			h.log.WithFields(logrus.Fields{
				"ID":      reqTask.ID,
				"Date":    reqTask.Date,
				"Title":   reqTask.Title,
				"Comment": reqTask.Comment,
				"Repeat":  reqTask.Repeat,
			}).Errorf("Ошибка при вычислении следующей даты: %v", err)
			h.SendResponse(w, models.APIResponse{
				StatusCode: http.StatusBadRequest,
				Body: models.Body{
					Error: "Ошибка при вычислении следующей даты",
				},
			})

			return
		}
		h.log.WithFields(logrus.Fields{"nextdate": nextdate}).Info("Следующая дата вычислена")

		h.log.WithFields(logrus.Fields{
			"ID":               reqTask.ID,
			"Date":             reqTask.Date,
			"Title":            reqTask.Title,
			"Comment":          reqTask.Comment,
			"Repeat":           reqTask.Repeat,
			"date.Before(now)": date.After(now),
			"date":             date,
			"now":              now,
		}).Info("date.Before(now)")

		//date.Before(now)
		//date.After(now) {
		if reqTask.Date < now.Format("20060102") {
			reqTask.Date = nextdate
		}

		//	}

	}

	//Если дата меньше сегодняшнего числа
	//если правило повторения не указано или равно пустой строке, подставляется сегодняшнее число;
	if date.Before(now) && reqTask.Repeat == "" {
		reqTask.Date = now.Format(scheduler.Format)
		h.log.Infof("Date is before now: %s, using current date", reqTask.Date)
	}

	// if nextdate != "" {

	// }

	id, err := h.storage.AddTask(ctx, reqTask)
	if err != nil {
		h.log.Errorf("Ошибка при создании задачи: %v", err)
		h.SendResponse(w, models.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Body: models.Body{
				Error:   "Internal server error",
				Message: "Ошибка при создании задачи",
			},
		})

		return
	}
	h.log.WithFields(logrus.Fields{"ID": reqTask.ID}).Info("Задача успешно создана")
	if id != 0 {
		h.SendResponse(w, models.APIResponse{
			StatusCode: http.StatusCreated,
			Body: models.NewTask{
				ID: id,
			},
		})
		return
	}
	h.SendResponse(w, models.APIResponse{
		StatusCode: http.StatusInternalServerError,
		Body: models.Body{
			Error:   "Internal server error",
			Message: "Ошибка при создании задачи",
		},
	})

}
