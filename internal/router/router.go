package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"hetic/tech-race/internal/handlers"
	"hetic/tech-race/internal/services"
	"log"
	"net/http"
	"text/template"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	templ := template.Must(template.ParseFiles("views/index.html"))
	err := templ.Execute(w, nil)
	if err != nil {
		println(err)
		return
	}
}

func SetupRouter(sessionService *services.SessionService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	sessionHandler := handlers.NewSessionHandler(sessionService)

	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Serve the home page
	r.Get("/", serveHome)

	r.Get("/sessions", sessionHandler.GetAll())
	r.Get("/sessions/start/{is_autopilot}", sessionHandler.Start())
	r.Get("/sessions/stop", sessionHandler.Stop())
	r.Get("/sessions/info", sessionHandler.GetAllSessionInfo())
	r.Get("/export-video", sessionHandler.ExportVideoToCloudinary())

	return r
}
