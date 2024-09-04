package config

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

type Config struct {
	DatabaseUrl         string
	ServerAddr          string
	CloudinaryID        string
	CloudinaryUrl       string
	CloudinaryUploadUrl string
}

func LoadFile() *Config {
	databaseURL := getEnv("DATABASE_URL", "postgres://root:password@localhost:5432/tech_race?sslmode=disable")
	serverAddr := getEnv("SERVER_ADDR", ":9000")
	cloudinaryUrl := getEnv("CLOUDINARY_URL", "")
	cloudinaryID := getEnv("CLOUDINARY_ID", "")
	cloudinaryUploadUrl := getEnv("CLOUDINARY_UPLOAD_URL", "http://localhost:8083/upload-video")

	config := &Config{
		DatabaseUrl:         databaseURL,
		ServerAddr:          serverAddr,
		CloudinaryID:        cloudinaryID,
		CloudinaryUrl:       cloudinaryUrl,
		CloudinaryUploadUrl: cloudinaryUploadUrl,
	}

	fmt.Printf("Loaded configuration : %s\n", config)

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
