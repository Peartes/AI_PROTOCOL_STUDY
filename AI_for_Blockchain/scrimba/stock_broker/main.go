package main

import (
	openai "github.com/openai/openai-go"
	openoptions "github.com/openai/openai-go/option"
)

var client openai.Client

func init() {
	client = *openai.NewClient(openoptions.WithAPIKey(OpenAIKey), openoptions.WithOrganization(OpenAIOrganization), openoptions.WithProject(OpenAIProject))
}

func main() {

}
