package entity

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

type ContextEntity struct {
	Cld       *cloudinary.Cloudinary
	Ctx       context.Context
	CloudName string
}

func NewContextEntity() *ContextEntity {
	errdot := godotenv.Load()
	if errdot != nil {
		log.Fatal("Vous avez besoin d'un fichier .env et d'y mettre l'accés cloudinary", errdot)
	}

	cloudName := os.Getenv("CLOUDINARY_ID")

	cld, err := cloudinary.New()
	if err != nil {
		log.Fatalf("Vous avez besoin de credentials cloudinary valides, %v", err)
	}
	cld.Config.URL.Secure = true
	ctx := context.Background()

	return &ContextEntity{Cld: cld, Ctx: ctx, CloudName: cloudName}
}

func (e *ContextEntity) UploadImage(imageURL, publicID string) (*uploader.UploadResult, error) {
	return e.Cld.Upload.Upload(e.Ctx, imageURL, uploader.UploadParams{
		PublicID:       publicID,
		UniqueFilename: api.Bool(false),
		Overwrite:      api.Bool(true),
	})
}

func (e *ContextEntity) UploadVideo(videoURL, publicID string) (*uploader.UploadResult, error) {
	return e.Cld.Upload.Upload(e.Ctx, videoURL, uploader.UploadParams{
		PublicID:       publicID,
		UniqueFilename: api.Bool(false),
		Overwrite:      api.Bool(true),
		ResourceType:   "video",
	})
}
