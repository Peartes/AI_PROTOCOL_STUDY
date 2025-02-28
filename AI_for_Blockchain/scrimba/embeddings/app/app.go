package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/openai/openai-go"
	openoptions "github.com/openai/openai-go/option"
	"github.com/peartes/scrimba/embeddings/config"
	supabase "github.com/supabase-community/supabase-go"
)

var client openai.Client
var ctx context.Context
var supaClient *supabase.Client

func init() {
	client = *openai.NewClient(openoptions.WithAPIKey(config.GetOpenAIKey()), openoptions.WithOrganization(config.GetOpenAIOrganization()), openoptions.WithProject(config.GetOpenAIProject()))
	ctx = context.Background()
	var err error
	supaClient, err = supabase.NewClient(config.GetSupaBaseUrl(), config.GetSupaBaseAPIKey(), nil)
	if err != nil {
		panic(fmt.Errorf("error creating supabase client %v", err))
	}
}

type vector struct {
	Content   string    `json:"content"`
	Embedding []float64 `json:"embedding"`
}

func RunApp() error {
	podcastData := config.GetPodcastMockup()

	var vectors []vector
	var wg sync.WaitGroup

	resultCh := make(chan vector, len(podcastData))

	for _, podcast := range podcastData {
		podcast := podcast
		wg.Add(1)
		go func(podcast string) {
			defer wg.Done()
			res, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
				Input: openai.Raw[openai.EmbeddingNewParamsInputUnion](podcast),
				Model: openai.Raw[openai.EmbeddingModel](openai.EmbeddingModelTextEmbeddingAda002),
			})
			if err != nil {
				println(fmt.Errorf("error calling embeddings endpoint for content %s %v ", podcast, err))
				resultCh <- vector{Content: podcast, Embedding: []float64{}}
				return
			}
			resultCh <- vector{Content: podcast, Embedding: res.Data[0].Embedding}
		}(podcast)
	}
	wg.Wait()
	close(resultCh)

	for v := range resultCh {
		vectors = append(vectors, v)
	}

	res := supaClient.From("documents").Insert(vectors, true, "", "", "")
	_, _, err := res.Execute()
	if err != nil {
		return fmt.Errorf("error inserting vectors into supabase %v", err)
	}
	return nil
}
