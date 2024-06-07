package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"mime/multipart"

	"github.com/cloudinarace/entity"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

type File struct {
	File multipart.File `json:"file,omitempty" validate:"required"`
}

type Url struct {
	Url string `json:"url,omitempty" validate:"required"`
}

type Response struct {
	StatusCode int                    `json:"statusCode"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data"`
}

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

func UploadVideoHandlerGin(entity *entity.ContextEntity) gin.HandlerFunc {
	return func(c *gin.Context) {

		videoURL := c.Query("url")
		videoPublicID := c.Query("id")
		uploadResult, err := entity.UploadVideo(videoURL, videoPublicID)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				Response{
					StatusCode: http.StatusInternalServerError,
					Message:    "error",
					Data:       map[string]interface{}{"data": "Erreur dans le téléchargement de la vidéo", "error": err},
				})
			return
		}

		c.JSON(
			http.StatusOK,
			Response{
				StatusCode: http.StatusOK,
				Message:    "success",
				Data:       map[string]interface{}{"data": uploadResult},
			})
	}
}

// ----------------------------------------------------------------
