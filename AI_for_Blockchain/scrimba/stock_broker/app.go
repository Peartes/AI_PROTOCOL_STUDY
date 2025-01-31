package main

import (
	"context"
	"fmt"

	openai "github.com/openai/openai-go"
	openoptions "github.com/openai/openai-go/option"
	"github.com/peartes/scrimba/stock_broker/config"
)

var client openai.Client
var ctx context.Context

func init() {
	client = *openai.NewClient(openoptions.WithAPIKey(config.GetOpenAIKey()), openoptions.WithOrganization(config.GetOpenAIOrganization()), openoptions.WithProject(config.GetOpenAIProject()))
	ctx = context.Background()
}

func chatMessage(role openai.ChatCompletionMessageParamRole, content string) openai.ChatCompletionMessageParam {
	return openai.ChatCompletionMessageParam{
		Role:    openai.Raw[openai.ChatCompletionMessageParamRole](role),
		Content: openai.Raw[interface{}](content),
	}
}

func RunApp() error {
	// get stock data
	tickers := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA"}
	channel, err := FetchStockData(tickers)
	if err != nil {
		return fmt.Errorf("error while fetching stock data %w ", err)
	}
	var stockData []struct {
		Symbol string
		PolygonAPIStockInfo
	}
	for data := <-channel; data != nil; data = <-channel {
		stockData = append(stockData, struct {
			Symbol string
			PolygonAPIStockInfo
		}{data.Ticker, data.Results[0]})
	}
	res, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.String(openai.ChatModelGPT3_5Turbo),
		Messages: openai.Raw[[]openai.ChatCompletionMessageParamUnion]([]openai.ChatCompletionMessageParamUnion{
			chatMessage(openai.ChatCompletionMessageParamRoleSystem, "You are a helpful general knowledge expert."),
			chatMessage(openai.ChatCompletionMessageParamRoleUser, fmt.Sprintf("Here is a summary of some stock info \n %v. \n"+
				"Use the data to analyze the stocks performance over the time period. "+
				"Be concise", stockData)),
		}),
	})
	if err != nil {
		fmt.Println(fmt.Errorf("error while asking model for chat completion %w ", err))
	} else {
		fmt.Printf("%s: %s \n", res.Choices[0].Message.Role, res.Choices[0].Message.Content)
	}
	return nil
}
