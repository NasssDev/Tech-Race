package handlers

import (
	"hetic/tech-race/internal/services"
	"net/http"
)

type VideosHandler struct {
	cloudinaryService *services.CloudinaryService
}

func NewVideosHandler(cloudinaryService *services.CloudinaryService) *VideosHandler {
	return &VideosHandler{cloudinaryService: cloudinaryService}
}

func (h *VideosHandler) GetAllVideos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// write a response to the client
		w.Write([]byte("Videos urls are ready to be queried !\n"))
	}
}
