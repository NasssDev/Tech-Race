package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"hetic/tech-race/internal/services"
	"hetic/tech-race/pkg/util"
	"net/http"
	"runtime"
	"strconv"
)

type SessionHandler struct {
	sessionService *services.SessionService
	uploadService  *services.UploadService
}

func NewSessionHandler(sessionService *services.SessionService, uploadService *services.UploadService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService, uploadService: uploadService}
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
			util.RenderJson(w, http.StatusOK, map[string]string{"status": "error", "message": "There is an error when parsing autopilot", "autopilot": ""})
			//http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		isActive, err := h.sessionService.IsSessionActive()
		if err != nil {
			util.RenderJson(w, http.StatusOK, map[string]string{"status": "error", "message": "There is an error when starting the car", "autopilot": isAutopilotStr})
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if isActive {
			util.RenderJson(w, http.StatusBadRequest, map[string]string{"status": "success", "message": "Autopilot mode is on", "autopilot": isAutopilotStr})
			//http.Error(w, "A session is already active", http.StatusBadRequest)
			return
		}

		err = h.sessionService.Start(isAutopilot)
		if err != nil {
			util.RenderJson(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": "Autopilot starting encounter an error", "autopilot": ""})
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		videoservice := services.NewVideoService(runtime.GOOS)

		//_, err = videoservice.StartRecording(h.sessionService)
		//if err != nil {
		//	util.RenderJson(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": "Car video cannot record", "autopilot": isAutopilotStr, "recording": ""})
		//	//w.Write([]byte("Session started with no recording \n"))
		//	fmt.Println("Error starting recording while session started:", err)
		//	//http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
		go func() {
			_, err := videoservice.StartRecording(h.sessionService)
			if err != nil {
				fmt.Println("Error starting recording while session started:", err)
			}
		}()

		fmt.Println("session started")
		// send json response with success message
		util.RenderJson(w, http.StatusOK, map[string]string{"status": "success", "message": "Autopilot mode is on", "autopilot": isAutopilotStr, "recording": "active"})
	}
}

func (h *SessionHandler) Stop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.sessionService.Stop()
		if err != nil {
			util.RenderJson(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": "There is an error in stopping the car", "autopilot": "", "recording": ""})
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		videoservice := services.NewVideoService(runtime.GOOS)
		videoservice.IsRecording = false
		fmt.Println("session stopped")
		//w.Write([]byte("Session stopped\n"))
		util.RenderJson(w, http.StatusOK, map[string]string{"status": "success", "message": "Car have been stopped", "autopilot": "", "recording": ""})
	}
}
func (h *SessionHandler) GetAllSessionInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionInfos, err := h.sessionService.GetAllSessionInfo()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			//TODO: que faire ici (on le stop ou non ?)
			//return
		}

		util.RenderJson(w, http.StatusOK, sessionInfos)
	}
}
