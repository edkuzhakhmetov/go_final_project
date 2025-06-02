package api

import (
	"encoding/json"
	"net/http"

	"github.com/edkuzhakhmetov/go_final_project/internal/models"
	"github.com/sirupsen/logrus"
)

func (h *Handler) SendResponse(w http.ResponseWriter, response models.APIResponse) {
	h.log.WithFields(logrus.Fields{
		"StatusCode": response.StatusCode,
		"Body":       response.Body,
		"Headers":    response.Headers,
	}).Info("Starting SendResponse action")

	if len(response.Headers) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	} else {
		for k, v := range response.Headers {
			w.Header().Set(k, v)
		}
	}
	w.WriteHeader(response.StatusCode)

	if response.Body != nil {
		body, err := json.Marshal(response.Body)

		if response.StatusCode >= 400 {
			h.log.WithFields(logrus.Fields{
				"StatusCode": response.StatusCode,
				"Body":       response.Body,
				"Headers":    response.Headers,
			}).Error("Error response")
		}

		w.Write(body)
		if err != nil {
			http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
			h.log.Errorf("failed to marshal response body: %w", err)
		}
	}

}
