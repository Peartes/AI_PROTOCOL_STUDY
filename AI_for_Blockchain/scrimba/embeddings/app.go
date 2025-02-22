package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/internal/param"
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
		Input: param.Field[openai.EmbeddingNewParamsInputUnion]("dsvsf"),
		Model: param.Field[openai.EmbeddingModel](openai.EmbeddingModelTextEmbeddingAda002),
	})
	if err != nil {
		return fmt.Errorf("error calling embeddings endpoint %v ", err)
	}
	return nil
}
