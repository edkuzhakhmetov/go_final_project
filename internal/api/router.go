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
	r.Get("/api/nextdate", handler.apiDelNextDay)
	r.Post("/api/task", handler.apiPostTask)
	r.Get("/api/tasks", handler.apiGetTasks)
	r.Get("/api/task", handler.apiGetTask)
	r.Put("/api/task", handler.apiPutTask)
	r.Delete("/api/task", handler.apiDelTask)
	r.Post("/api/task/done", handler.apiPostTaskDone)

	return &Router{
		Mux:     r,
		Handler: handler,
	}
}
