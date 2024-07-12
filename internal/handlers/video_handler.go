package handlers

import (
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
		var cloudinaryPackageUrl = "http://localhost:8045/upload-video"
		//url=&id=2024-07-11T16:29:13
		resp := services.UploadVideoToCloudinary(cloudinaryPackageUrl, "../../../tmp/video/2024-07-11T16:29:13.mp4", "2024-07-11T16:29:13")

		//timestamp := time.Now()
		util.RenderJson(w, http.StatusOK, resp)

	}
}
