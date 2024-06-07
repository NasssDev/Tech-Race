package handlers

import (
	"hetic/tech-race/internal/services"
	"net/http"
)

type SessionHandler struct {
	sessionService *services.SessionService
}

func NewSessionHandler(sessionService *services.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

func (h *SessionHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// write a response to the client
		w.Write([]byte("Tech race is ready to GOOO !\n"))
	}
}
