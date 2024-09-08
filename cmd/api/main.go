package main

import (
	"hetic/tech-race/internal/config"
	"hetic/tech-race/internal/database"
	"hetic/tech-race/internal/router"
	"hetic/tech-race/internal/services"
	"log"
	"mime"
	"net/http"
)

func main() {
	cfg := config.LoadAppInfo()

	db := database.Connect(cfg.DatabaseUrl)

	dbWrapper := database.NewDatabase(db)

	sessionService := services.NewSessionService(dbWrapper)
	uploadService := services.NewUploadService(dbWrapper)

	err := mime.AddExtensionType(".css", "text/css")
	if err != nil {
		log.Fatal("Error adding MIME type for .css files:", err)
	}

	r := router.SetupRouter(sessionService, uploadService)

	log.Printf("Starting server on %s", cfg.ServerAddr)
	log.Fatal(http.ListenAndServe(cfg.ServerAddr, r))
}
