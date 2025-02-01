package main

import openai "github.com/openai/openai-go"

func ChatMessage(role openai.ChatCompletionMessageParamRole, content string) openai.ChatCompletionMessageParam {
	return openai.ChatCompletionMessageParam{
		Role:    openai.Raw[openai.ChatCompletionMessageParamRole](role),
		Content: openai.Raw[interface{}](content),
	}
}
