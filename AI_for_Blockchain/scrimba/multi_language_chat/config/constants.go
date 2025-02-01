package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetOpenAIKey() string {
	return os.Getenv("OPENAI_API_KEY")
}

func GetOpenAIOrganization() string {
	return os.Getenv("OPENAI_ORGANIZATION")
}

func GetOpenAIProject() string {
	return os.Getenv("OPENAI_PROJECT")
}
