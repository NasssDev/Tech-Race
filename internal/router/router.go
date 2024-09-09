package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"hetic/tech-race/internal/handlers"
	"hetic/tech-race/internal/services"
	"net/http"
)

func SetupRouter(sessionService *services.SessionService, uploadService *services.UploadService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	sessionHandler := handlers.NewSessionHandler(sessionService, uploadService)

	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Serve the html pages
	r.Get("/", handlers.ServeHome)
	r.Get("/docs", handlers.ServeDocs)
	r.Get("/tarifs", handlers.ServeTarifs)

	r.Get("/sessions", sessionHandler.GetAll())
	r.Get("/sessions/start/{is_autopilot}", sessionHandler.Start())
	r.Get("/sessions/stop", sessionHandler.Stop())
	r.Get("/sessions/info", sessionHandler.GetAllSessionInfo())

	return r
}
