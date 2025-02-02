package ai

import (
	"context"
	"fmt"

	openai "github.com/openai/openai-go"
	openoptions "github.com/openai/openai-go/option"
	"github.com/peartes/scrimba/pollyglot/config"
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

func Translate(text, targetLanguage string) (string, error) {
	res, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.String(openai.ChatModelGPT3_5Turbo),
		Messages: openai.Raw[[]openai.ChatCompletionMessageParamUnion]([]openai.ChatCompletionMessageParamUnion{
			chatMessage(openai.ChatCompletionMessageParamRoleSystem, fmt.Sprintf("You are strictly a chatbot language translator. You translate text to %s. You don't embellish words but do a word for word translation or phrase translation when that makes more sense. If the language to translate is the same as the language to translate into, return the text as is. If you cannot translate a text, return as is", targetLanguage)),
			chatMessage(openai.ChatCompletionMessageParamRoleUser, text),
		}),
	})
	if err != nil {
		return "", fmt.Errorf("error while asking model for chat completion %w ", err)
	}
	return res.Choices[0].Message.Content, nil
}
