package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/edkuzhakhmetov/go_final_project/internal/models"
	"github.com/edkuzhakhmetov/go_final_project/internal/storage"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	log     logrus.FieldLogger
	storage *storage.Storage
}

func NewHandler(logger logrus.FieldLogger, storage *storage.Storage) *Handler {
	return &Handler{
		log:     logger,
		storage: storage,
	}
}

func (h *Handler) getTaskFromBody(r *http.Request) (models.Task, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.log.Infof("Не удалось распарсить тело запроса: %v", err)
		return models.Task{}, err
	}
	var task models.Task
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		h.log.Infof("Не удалось распарсить тело запроса: %v", err)
		return models.Task{}, err
	}

	h.log.WithFields(logrus.Fields{
		"ID":      task.ID,
		"Date":    task.Date,
		"Title":   task.Title,
		"Comment": task.Comment,
		"Repeat":  task.Repeat,
	}).Info("Получен запрос c задачей")

	return task, nil
}
