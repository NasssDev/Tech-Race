package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"hetic/tech-race/internal/handlers"
	"hetic/tech-race/internal/services"
)

func SetupRouter(sessionService *services.SessionService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	sessionHandler := handlers.NewSessionHandler(sessionService)

	r.Get("/sessions", sessionHandler.GetAll())
	r.Get("/sessions/start/{is_autopilot}", sessionHandler.Start())
	r.Get("/sessions/stop", sessionHandler.Stop())
	r.Get("/sessions/info", sessionHandler.GetAllSessionInfo())
	r.Get("/export-video", sessionHandler.ExportVideoToCloudinary())

	return r
}
