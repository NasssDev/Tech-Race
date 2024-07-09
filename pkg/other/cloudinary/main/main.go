package main

import (
	"fmt"
	"github.com/cloudinarace/entity"
	"github.com/cloudinarace/handler"
	"github.com/joho/godotenv"
	"log"
	"mime"
	"net/http"
	"os"
	"text/template"
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
	errdot := godotenv.Load()
	if errdot != nil {
		log.Fatal("Error loading .env file:", errdot)
	}

	// Add MIME type for .css files
	err := mime.AddExtensionType(".css", "text/css")
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

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/upload", handler.UploadImageHandler(entity.NewContextEntity()))
	http.HandleFunc("/info", handler.GetAssetInfoHandler(entity.NewContextEntity()))
	http.HandleFunc("/transform", handler.TransformImageHandler(entity.NewContextEntity()))
	http.HandleFunc("/display", handler.DisplayImageHandler(entity.NewContextEntity()))
	http.HandleFunc("/upload-video", handler.UploadVideoHandler(entity.NewContextEntity()))
	http.HandleFunc("/display-video", handler.DisplayVideoHandler(entity.NewContextEntity()))

	fmt.Println("Starting server on server " + os.Getenv("PORT") + "...")
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
