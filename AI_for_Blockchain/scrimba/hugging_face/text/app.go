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

func TextClassification() error {
	res, err := client.TextClassification(ctx, &hf.TextClassificationRequest{
		Inputs: "Today is a good day",
	})
	if err != nil {
		return fmt.Errorf("error classifying text %w", err)
	}
	fmt.Printf("Classification is %s ", res[0][0].Label)
	return nil
}

func TextTranslation(srcText, srcLanguage string) error {
	res, err := client.Translation(ctx, &hf.TranslationRequest{
		Inputs: []string{srcText},
		Model:  "facebook/mbart-large-50-many-to-many-mmt",
	})
	if err != nil {
		return fmt.Errorf("error classifying text %w", err)
	}
	fmt.Printf("translation is %s ", res[0].TranslationText)
	return nil
}

func RunApp() error {
	if err := TextTranslation("Меня зовут Вольфганг и я живу в Берлине", "en_XX"); err != nil {
		return fmt.Errorf("error running app: %w", err)
	}
	return nil
}
