package llm

import (
	"context"
	"log"

	ark "github.com/sashabaranov/go-openai"
)

type OpenAICompatibleClient struct {
	Model  string
	Client *ark.Client
}

func (o *OpenAICompatibleClient) Generate(prompt, systemPrompt string) string {
	log.Println("正在调用大语言模型...")

	message := []ark.ChatCompletionMessage{
		{Role: ark.ChatMessageRoleSystem, Content: systemPrompt},
		{Role: ark.ChatMessageRoleUser, Content: prompt},
	}

	response, err := o.Client.CreateChatCompletion(
		context.Background(),
		ark.ChatCompletionRequest{
			Model:    o.Model,
			Messages: message,
			Stream:   false,
		},
	)
	if err != nil {
		log.Println(err)
		return "调用大模型出错" + err.Error()
	}

	answer := response.Choices[0].Message.Content
	log.Println("调用大模型响应成功")
	return answer
}
