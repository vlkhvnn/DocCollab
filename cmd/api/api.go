package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	ws "github.com/vlkhvnn/DocCollab/internal/ws/hub"
)

type application struct {
	config config
}

type config struct {
	addr string
	hub  *ws.Hub
}

func mount(app *application) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Get("/ws", app.serveWs)
	})
	return r
}

// homeHandler provides a simple home page.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("WebSocket server is running. Connect to /ws"))
}
