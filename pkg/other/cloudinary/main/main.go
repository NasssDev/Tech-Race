package main

import (
	"flag"
	"fmt"
	"github.com/cloudinarace/config"
	"github.com/cloudinarace/entity"
	"github.com/cloudinarace/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"mime"
	"net/http"
	"os"
)

func main() {
	port := flag.String("port", config.GetEnv("PORT", "9080"), "Port to run the server on")
	flag.Parse()

	// Set the parsed port in the environment variable
	err := os.Setenv("PORT", *port)
	if err != nil {
		log.Fatalf("Error setting environment variable: %v", err)
	}

	// Load environment variables from .env file
	errdot := godotenv.Load()
	if errdot != nil {
		log.Fatal("Error loading .env file:", errdot)
	}

	// Add MIME types for .css and .js files
	err = mime.AddExtensionType(".css", "text/css")
	if err != nil {
		log.Fatalf("Error adding MIME type for .css: %v", err)
	}

	err = mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		log.Fatalf("Error adding MIME type for .js: %v", err)
	}

	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("views/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Define other routes using Gin handlers
	r.GET("/upload-video", handler.UploadVideoHandlerGin(entity.NewContextEntity()))
	r.POST("/upload", gin.WrapF(handler.UploadImageHandler(entity.NewContextEntity())))
	r.GET("/info", gin.WrapF(handler.GetAssetInfoHandler(entity.NewContextEntity())))
	r.POST("/transform", gin.WrapF(handler.TransformImageHandler(entity.NewContextEntity())))
	r.GET("/display", gin.WrapF(handler.DisplayImageHandler(entity.NewContextEntity())))
	r.GET("/display-video", gin.WrapF(handler.DisplayVideoHandler(entity.NewContextEntity())))

	// Start the Gin server on the specified port
	portEnv := os.Getenv("PORT")
	fmt.Println("Starting server on port " + portEnv + "...")
	if err := r.Run(":" + portEnv); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
