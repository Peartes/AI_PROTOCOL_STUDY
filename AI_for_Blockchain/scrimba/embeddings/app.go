package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	openoptions "github.com/openai/openai-go/option"
	"github.com/peartes/scrimba/embeddings/config"
)

var client openai.Client
var ctx context.Context

func init() {
	client = *openai.NewClient(openoptions.WithAPIKey(config.GetOpenAIKey()), openoptions.WithOrganization(config.GetOpenAIOrganization()), openoptions.WithProject(config.GetOpenAIProject()))
	ctx = context.Background()
}

func RunApp() error {
	res, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.Raw[openai.EmbeddingNewParamsInputUnion]("Hello World!"),
		Model: openai.Raw[openai.EmbeddingModel](openai.EmbeddingModelTextEmbeddingAda002),
	})
	fmt.Printf("res: %v", res)
	if err != nil {
		return fmt.Errorf("error calling embeddings endpoint %v ", err)
	}
	return nil
}
