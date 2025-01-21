package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

var (
	// OpenAIKey is the key to access the OpenAI API
	OpenAIKey = os.Getenv("OPENAI_API_KEY")
	// OpenAIOrganization is the organization to access the OpenAI API
	OpenAIOrganization = os.Getenv("OPENAI_ORGANIZATION")
	// OpenAIProject is the project to access the OpenAI API
	OpenAIProject = os.Getenv("OPENAI_PROJECT")
)
