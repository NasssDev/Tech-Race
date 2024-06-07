package handler

import (
	"fmt"
	"github.com/cloudinarace/entity"
	"net/http"
)

func TransformImageHandler(entity *entity.ContextEntity) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		qsImg, err := entity.Cld.Image("quickstart_butterfly")
		if err != nil {
			http.Error(w, "Error creating image object", http.StatusInternalServerError)
			return
		}

		qsImg.Transformation = "r_max/e_sepia"
		newURL, err := qsImg.String()
		if err != nil {
			http.Error(w, "Error transforming image", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Transformed image URL: %s\n", newURL)
	}
}
