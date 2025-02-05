package main

import (
	"context"
	"fmt"

	hf "github.com/hupe1980/go-huggingface"
	"github.com/peartes/scrimba/hugging_face/text/config"
)

var client *hf.InferenceClient
var ctx context.Context

func init() {
	client = hf.NewInferenceClient(config.GetHuggingFaceToken())
	ctx = context.Background()
}

func TextGeneration() error {
	res, err := client.TextGeneration(ctx, &hf.TextGenerationRequest{
		Model:  "HuggingFaceH4/zephyr-7b-beta",
		Inputs: "what is the meaning of life?",
	})
	if err != nil {
		return fmt.Errorf("error generating text: %w", err)
	}
	fmt.Printf("Generated text: %s\n", res[0].GeneratedText)
	return nil
}

func RunApp() error {
	if err := TextGeneration(); err != nil {
		return fmt.Errorf("error running app: %w", err)
	}
	return nil
}
