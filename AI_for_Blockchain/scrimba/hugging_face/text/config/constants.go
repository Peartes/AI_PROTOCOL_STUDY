package config

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}
}

func GetHuggingFaceToken() string {
	return os.Getenv("HUGGING_FACE_API_KEY")
}
