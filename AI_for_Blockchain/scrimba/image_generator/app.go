package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	openoptions "github.com/openai/openai-go/option"
	"github.com/peartes/scrimba/image_generator/config"
)

var client openai.Client
var ctx context.Context

func init() {
	client = *openai.NewClient(openoptions.WithAPIKey(config.GetOpenAIKey()), openoptions.WithOrganization(config.GetOpenAIOrganization()), openoptions.WithProject(config.GetOpenAIProject()))
	ctx = context.Background()
}

func RunApp(prompt string) error {
	res, err := client.Images.Generate(ctx, openai.ImageGenerateParams{
		Prompt:         openai.String(prompt),
		Model:          openai.String(openai.ImageModelDallE2),
		N:              openai.Int(1),
		ResponseFormat: openai.Raw[openai.ImageGenerateParamsResponseFormat](openai.ImageGenerateParamsResponseFormatB64JSON),
		Quality:        openai.Raw[openai.ImageGenerateParamsQuality](openai.ImageGenerateParamsQualityStandard),
		Size:           openai.Raw[openai.ImageGenerateParamsSize](openai.ImageGenerateParamsSize256x256),
		Style:          openai.Raw[openai.ImageGenerateParamsStyle](openai.ImageGenerateParamsStyleNatural),
	})
	if err != nil {
		return fmt.Errorf("error while asking model for image generation %w ", err)
	}
	fmt.Printf("Generate Image URL: %s \n", res.Data[0].URL)
	return nil
}
