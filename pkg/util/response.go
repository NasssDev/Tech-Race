package util

import (
	"encoding/json"
	"net/http"
	"text/template"
)

func RenderJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ServeHTML(HtmlUrl string) *template.Template {

	templ := template.Must(template.ParseFiles(HtmlUrl))
	return templ
}
