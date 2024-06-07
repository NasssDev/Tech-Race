package main

import (
	"fmt"
	"github.com/cloudinarace/entity"
	"github.com/cloudinarace/handler"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	errdot := godotenv.Load()
	if errdot != nil {
		log.Fatal("Error loading .env file:", errdot)
	}

	http.HandleFunc("/upload", handler.UploadImageHandler(entity.NewContextEntity()))
	http.HandleFunc("/info", handler.GetAssetInfoHandler(entity.NewContextEntity()))
	http.HandleFunc("/transform", handler.TransformImageHandler(entity.NewContextEntity()))
	http.HandleFunc("/display", handler.DisplayImageHandler(entity.NewContextEntity()))
	http.HandleFunc("/upload-video", handler.UploadVideoHandler(entity.NewContextEntity()))
	http.HandleFunc("/display-video", handler.DisplayVideoHandler(entity.NewContextEntity()))

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
