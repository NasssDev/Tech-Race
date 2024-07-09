package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"hetic/tech-race/internal/services"
	"hetic/tech-race/pkg/util"
	"net/http"
	"strconv"
)

type SessionHandler struct {
	sessionService *services.SessionService
}

func NewSessionHandler(sessionService *services.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

func (h *SessionHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Tech race is ready to GOOO !\n"))
		sessions, err := h.sessionService.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(sessions)

		util.RenderJson(w, http.StatusOK, sessions)
	}
}
func (h *SessionHandler) Start() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAutopilotStr := chi.URLParam(r, "is_autopilot")
		isAutopilot, err := strconv.ParseBool(isAutopilotStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		isActive, err := h.sessionService.IsSessionActive()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if isActive {
			http.Error(w, "A session is already active", http.StatusBadRequest)
			return
		}

		err = h.sessionService.Start(isAutopilot)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Session started\n"))
	}
}

func (h *SessionHandler) Stop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.sessionService.Stop()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Session stopped\n"))
	}
}
func (h *SessionHandler) GetAllSessionInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionInfos, err := h.sessionService.GetAllSessionInfo()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		util.RenderJson(w, http.StatusOK, sessionInfos)
	}
}
