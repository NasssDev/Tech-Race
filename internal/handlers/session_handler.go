package handlers

import (
	"fmt"
	"hetic/tech-race/internal/services"
	"hetic/tech-race/pkg/util"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/go-chi/chi/v5"
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//currentSessionId, err := h.sessionService.GetCurrentSessionID()
		//if err != nil {
		//println("problem getting session id")
		//}

		//sessionId := currentSessionId.ID

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

		// Lance le relay du stream de l'esp32 permettant de gérer plusieurs connections simultanée sur le port 8082
		go func() {
			cmd := exec.Command("go", "run", "cmd/stream/main.go")
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error running stream:", err)
			}
		}()

		videoservice := services.NewVideoService(runtime.GOOS)
		//recordingData, err := videoservice.StartRecording(h.sessionService)
		_, err = videoservice.StartRecording(h.sessionService)
		if err != nil {
			fmt.Println("Error starting recording:", err)
		}

		//var cloudinaryPackageUrl = "http://localhost:8083/upload-video"
		//dirFromCloudinarace := "../../../tmp/video"

		//resp, err := services.UploadVideoToCloudinary(cloudinaryPackageUrl, filepath.Join(dirFromCloudinarace, recordingData.VideoName+".mp4"), recordingData.VideoName)
		//if err != nil {
		//	fmt.Println("Error uploading video:", err)
		//}
		//
		//videoPath := resp.Data.Data.URL
		//if videoPath != "" {
		//	err := h.uploadService.InsertVideo(videoPath)
		//	if err != nil {
		//		println("error in database insertion: ", err)
		//		http.Error(w, err.Error(), http.StatusInternalServerError)
		//	}
		//} else {
		//	http.Error(w, "L'url de la vidéo n'a pas été trouvé", http.StatusInternalServerError)
		//}

		// don't display this if StartRecording returned an error
		w.Write([]byte("Session started\n"))

		// send json response with success message
		util.RenderJson(w, http.StatusOK, map[string]string{"message": "Session started"})
	}
}

func (h *SessionHandler) Stop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.sessionService.Stop()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		videoservice := services.NewVideoService(runtime.GOOS)
		videoservice.IsRecording = false

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
