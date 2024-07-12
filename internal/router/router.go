package router

import (
	"hetic/tech-race/internal/handlers"
	"hetic/tech-race/internal/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(sessionService *services.SessionService, videoService *services.UploadService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	sessionHandler := handlers.NewSessionHandler(sessionService)
	videoHandler := handlers.NewVideoHandler(videoService)

	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Serve the home page
	r.Get("/", handlers.ServeHome)

	r.Get("/sessions", sessionHandler.GetAll())
	r.Get("/sessions/start/{is_autopilot}", sessionHandler.Start())
	r.Get("/sessions/stop", sessionHandler.Stop())
	r.Get("/sessions/info", sessionHandler.GetAllSessionInfo())
	r.Get("/export-video", videoHandler.ExportVideoToCloudinary())

	return r
}
