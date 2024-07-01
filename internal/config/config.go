package config

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

type Config struct {
	DatabaseUrl string
	ServerAddr  string
}

func LoadFile() *Config {
	databaseURL := getEnv("DATABASE_URL", "postgres://root:password@localhost:5432/tech_race?sslmode=disable")
	serverAddr := getEnv("SERVER_ADDR", ":9000")

	config := &Config{
		DatabaseUrl: databaseURL,
		ServerAddr:  serverAddr,
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
