package handlers

import (
	"hetic/tech-race/pkg/util"
	"log"
	"net/http"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {

	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	templ := util.ServeHTML("views/index.html")
	err := templ.Execute(w, nil)
	if err != nil {
		println(err)
		return
	}
}
