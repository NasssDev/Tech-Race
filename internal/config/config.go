package config

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

type AppInfo struct {
	DatabaseUrl         string
	ServerAddr          string
	CloudinaryID        string
	CloudinaryUrl       string
	CloudinaryUploadUrl string
}

type StreamInfo struct {
	Esp32Address   string
	Esp32Port      string
	RelayAddress   string
	StreamBoundary string
}

func LoadAppInfo() *AppInfo {
	databaseURL := getEnv("DATABASE_URL", "postgres://root:password@localhost:5432/tech_race?sslmode=disable")
	serverAddr := getEnv("SERVER_ADDR", ":9000")
	cloudinaryUrl := getEnv("CLOUDINARY_URL", "")
	cloudinaryID := getEnv("CLOUDINARY_ID", "")
	cloudinaryUploadUrl := getEnv("CLOUDINARY_UPLOAD_URL", "http://localhost:8083/upload-video")

	config := &AppInfo{
		DatabaseUrl:         databaseURL,
		ServerAddr:          serverAddr,
		CloudinaryID:        cloudinaryID,
		CloudinaryUrl:       cloudinaryUrl,
		CloudinaryUploadUrl: cloudinaryUploadUrl,
	}

	fmt.Printf("Loaded configuration : %s\n", config)

	return config
}

func LoadStreamInfo() *StreamInfo {
	esp32Address := getEnv("ESP32_ADDRESS", "192.168.16.10:7000")
	esp32Port := getEnv("ESP32_PORT", "7000")
	relayAddress := getEnv("RELAY_ADDRESS", ":8080")
	streamBoundary := getEnv("STREAM_BOUNDARY", "--123456789000000000000987654321")

	config := &StreamInfo{
		Esp32Address:   esp32Address,
		Esp32Port:      esp32Port,
		RelayAddress:   relayAddress,
		StreamBoundary: streamBoundary,
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
