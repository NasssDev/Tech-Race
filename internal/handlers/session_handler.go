package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"hetic/tech-race/internal/services"
	"hetic/tech-race/pkg/util"
	"net/http"
	"path/filepath"
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
			util.RenderJson(w, http.StatusBadRequest, map[string]string{"status": "error", "message": "Invalid autopilot parameter", "autopilot": ""})
			return
		}

		isActive, err := h.sessionService.IsSessionActive()
		if err != nil {
			util.RenderJson(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": "Error checking session status", "autopilot": ""})
			return
		}

		if isActive {
			util.RenderJson(w, http.StatusBadRequest, map[string]string{"status": "error", "message": "A session is already active", "autopilot": isAutopilotStr})
			return
		}

		err = h.sessionService.Start(isAutopilot)
		if err != nil {
			util.RenderJson(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": "Error starting session", "autopilot": ""})
			return
		}

		// Envoyer la réponse HTTP immédiatement
		util.RenderJson(w, http.StatusOK, map[string]string{"status": "success", "message": "Session started", "autopilot": isAutopilotStr, "recording": "starting"})

		// Démarrer l'enregistrement vidéo en arrière-plan
		go func() {
			videoservice := services.NewVideoService(runtime.GOOS)
			recordingData, err := videoservice.StartRecording(h.sessionService)
			if err != nil {
				fmt.Println("Error starting recording:", err)
				return
			}

			var cloudinaryPackageUrl = "http://localhost:8083/upload-video"
			dirFromCloudinarace := "../../../tmp/video"

			resp, err := services.UploadVideoToCloudinary(cloudinaryPackageUrl, filepath.Join(dirFromCloudinarace, recordingData.VideoName+".mp4"), recordingData.VideoName)
			if err != nil {
				fmt.Println("Error uploading video:", err)
				return
			}

			videoPath := resp.Data.Data.URL
			if videoPath != "" {
				err := h.uploadService.InsertVideo(videoPath)
				if err != nil {
					fmt.Println("Error inserting video in database:", err)
				}
			} else {
				fmt.Println("Video URL not found")
			}
		}()
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
