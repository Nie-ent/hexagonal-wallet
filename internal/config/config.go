package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(path string) {
	_ = godotenv.Load(path)
}

func GetEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
