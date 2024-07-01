package handler

import (
	"github.com/cloudinarace/entity"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"log"
	"net/http"
)

func GetAssetInfoHandler(entity *entity.ContextEntity) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		asset, err := entity.Cld.Admin.Asset(entity.Ctx, admin.AssetParams{PublicID: "quickstart_butterfly"})
		if err != nil {
			http.Error(w, "Failed to get asset details", http.StatusInternalServerError)
			return
		}

		// Print some basic information about the asset.
		log.Printf("Public ID: %v, URL: %v\n", asset.PublicID, asset.SecureURL)
	}
}
