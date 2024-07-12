package main

import (
	"flag"
	"log"
	"mime"
	"net/http"
	"os"
	"text/template"

	"github.com/cloudinarace/entity"
	"github.com/cloudinarace/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	templ := template.Must(template.ParseFiles("views/index.html"))
	err := templ.Execute(w, nil)
	if err != nil {
		println(err)
		return
	}
}

func main() {
	// Define and parse the port flag
	port := flag.String("port", "8090", "Port to run the server on")
	flag.Parse()

	// Set the parsed port in the environment variable
	err := os.Setenv("PORT", *port)
	if err != nil {
		return
	}

	// Load environment variables from .env file
	errdot := godotenv.Load()
	if errdot != nil {
		log.Fatal("Error loading .env file:", errdot)
	}

	// Add MIME type for .css files
	err = mime.AddExtensionType(".css", "text/css")
	if err != nil {
		println(err)
		return
	}

	// Add MIME type for .js files
	err = mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		println(err)
		return
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/upload", handler.UploadImageHandler(entity.NewContextEntity()))
	http.HandleFunc("/info", handler.GetAssetInfoHandler(entity.NewContextEntity()))
	http.HandleFunc("/transform", handler.TransformImageHandler(entity.NewContextEntity()))
	http.HandleFunc("/display", handler.DisplayImageHandler(entity.NewContextEntity()))

	http.HandleFunc("/display-video", handler.DisplayVideoHandler(entity.NewContextEntity()))

	// test : upload-video?url=../../../tmp/video/2024-07-11T16:29:13.mp4&id=2024-07-11T16:29:13
	r.GET("/upload-video", handler.UploadVideoHandlerGin(entity.NewContextEntity()))

	portEnv := os.Getenv("PORT")

	r.Run(":" + portEnv)

	// fmt.Println("Starting server on port " + portEnv + "...")
	// if err := http.ListenAndServe(":"+portEnv, nil); err != nil {
	// 	log.Fatalf("Server failed to start: %v", err)
	// }
}
