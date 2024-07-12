package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cloudinarace/entity"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
)

func GetAssetInfoHandler(entity *entity.ContextEntity) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		publicID := r.URL.Query().Get("id")
		if publicID == "" {
			http.Error(w, "Public ID parameter is required", http.StatusBadRequest)
			return
		}
		// Print the Public ID being used
		log.Printf("Attempting to retrieve asset with Public ID: %s\n", publicID)

		//asset, err := entity.Cld.Admin.Asset(entity.Ctx, admin.AssetParams{PublicID: "three_dogs"})
		asset, err := entity.Cld.Admin.Asset(entity.Ctx, admin.AssetParams{
			PublicID: publicID})
		if err != nil {
			http.Error(w, "Failed to get asset details", http.StatusInternalServerError)
			return
		}

		// Print some basic information about the asset.
		log.Printf("Public ID: %v, URL: %v\n", asset.PublicID, asset.SecureURL)

		response := map[string]interface{}{
			"message": "Successfully found the asset",
			"results": asset,
		}

		// Set content type to JSON
		w.Header().Set("Content-Type", "application/json")
		// Write JSON response
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
