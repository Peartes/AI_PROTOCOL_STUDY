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

func main() {
	res, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.String(openai.ChatModelGPT4_0613),
		Messages: openai.Raw[[]openai.ChatCompletionMessageParamUnion]([]openai.ChatCompletionMessageParamUnion{
			chatMessage(openai.ChatCompletionMessageParamRoleSystem, "You are a helpful general knowledge expert."),
			chatMessage(openai.ChatCompletionMessageParamRoleUser, "What is the capital of the United States?"),
		}),
	})
	if err != nil {
		fmt.Println(fmt.Errorf("error while asking model for chat completion %w ", err))
	} else {
		fmt.Printf("%s: %s \n", res.Choices[0].Message.Role, res.Choices[0].Message.Content)
	}
}
