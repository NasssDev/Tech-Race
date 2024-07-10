package handler

import (
	"fmt"
	"github.com/cloudinarace/entity"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"net/http"
	"strings"
)

func UploadImageHandler(entity *entity.ContextEntity) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		imageURL := r.URL.Query().Get("url")
		imagePublicID := r.URL.Query().Get("id")
		if imageURL == "" {
			http.Error(w, "URL parameter is required", http.StatusBadRequest)
			return
		}

		uploadResult, err := entity.Cld.Upload.Upload(
			entity.Ctx,
			imageURL,
			uploader.UploadParams{
				PublicID:       imagePublicID,
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
		// Extract the URL parameter from the request
		videoURL := r.URL.Query().Get("url")
		videoPublicID := r.URL.Query().Get("id")
		if videoURL == "" {
			http.Error(w, "Il manque des query string ?id= et ?url=", http.StatusBadRequest)
			return
		}

		if !strings.HasSuffix(videoURL, ".mp4") {
			http.Error(w, "Format obligatoire en .mp4", http.StatusBadRequest)
			return
		}

		uploadResult, err := entity.UploadVideo(videoURL, videoPublicID)
		if err != nil {
			println("Echec export de vidéos - erreur: ", err)
			http.Error(w, "Echec dans l' export des vidéos, statut: ", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Vos vidéos sont sur Cloudinary, bravo! URL: %s\n", uploadResult.SecureURL)
	}
}
