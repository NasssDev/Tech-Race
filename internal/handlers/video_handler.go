package handlers

import (
	"hetic/tech-race/internal/config"
	"hetic/tech-race/internal/services"
	"hetic/tech-race/pkg/util"
	"net/http"
)

type VideoHandler struct {
	videoService *services.UploadService
}

func NewVideoHandler(videoService *services.UploadService) *VideoHandler {
	return &VideoHandler{videoService: videoService}
}

func (h *VideoHandler) ExportVideoToCloudinary() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := config.LoadFile()
		var cloudinaryPackageUrl = cfg.CloudinaryUploadUrl

		resp, uploaderr := services.UploadVideoToCloudinary(cloudinaryPackageUrl, "../../../tmp/video/2024-07-11T16:29:13.mp4", "2024-07-11T16:29:13")
		if uploaderr != nil {
			http.Error(w, uploaderr.Error(), http.StatusInternalServerError)
			return
		}

		util.RenderJson(w, http.StatusOK, resp)

		videoPath := resp.Data.Data.URL
		if videoPath != "" {
			err := h.videoService.InsertVideo(videoPath)
			if err != nil {
				println("error in database insertion: ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "L'url de la vidéo n'a pas été trouvé", http.StatusInternalServerError)
		}

	}
}
