package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Router struct {
	Mux     *chi.Mux
	Handler *Handler
}

var webDir = "./web"

func NewRouter(handler *Handler) *Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, webDir+"/index.html")
	})

	fs := http.FileServer(http.Dir(webDir))
	r.Handle("/*", http.StripPrefix("/", fs))
	r.Get("/api/nextdate", handler.NextDayHandler)
	r.Post("/api/task", handler.PostTaskHandler)
	r.Get("/api/tasks", handler.getTasksHandler)
	r.Get("/api/task", handler.getTaskHandler)
	r.Put("/api/task", handler.putTaskHandler)
	r.Delete("/api/task", handler.delTaskHandler)
	r.Post("/api/task/done", handler.postTaskDoneHandler)

	return &Router{
		Mux:     r,
		Handler: handler,
	}
}
