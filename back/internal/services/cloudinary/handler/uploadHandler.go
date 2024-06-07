package handler

import (
	"fmt"
	"github.com/cloudinarace/entity"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"net/http"
)

func UploadImageHandler(entity *entity.ContextEntity) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uploadResult, err := entity.Cld.Upload.Upload(
			entity.Ctx,
			"https://cloudinary-devs.github.io/cld-docs-assets/assets/images/butterfly.jpeg",
			uploader.UploadParams{
				PublicID:       "quickstart_butterfly",
				UniqueFilename: api.Bool(false),
				Overwrite:      api.Bool(true)})
		if err != nil {
			http.Error(w, "Failed to upload file", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Image uploaded successfully! Delivery URL: %s\n", uploadResult.SecureURL)
	}
}

func UploadVideoHandler(entity *entity.ContextEntity) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uploadResult, err := entity.UploadVideo("https://res.cloudinary.com/demo/video/upload/wave.mp4", "quickstart_wave")
		if err != nil {
			http.Error(w, "Failed to upload video", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Video uploaded successfully! Delivery URL: %s\n", uploadResult.SecureURL)
	}
}
