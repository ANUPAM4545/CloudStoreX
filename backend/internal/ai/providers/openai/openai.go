package openai

import (
	"context"

	"github.com/cloudstorex/backend/internal/ai"
	sdk "github.com/sashabaranov/go-openai"
)

type provider struct {
	client *sdk.Client
}

// NewProvider creates a new OpenAI provider wrapper
func NewProvider(apiKey string) ai.AIProvider {
	return &provider{
		client: sdk.NewClient(apiKey),
	}
}

func (p *provider) GetProviderType() ai.ProviderType {
	return ai.ProviderOpenAI
}

func (p *provider) GenerateText(ctx context.Context, req ai.GenerateTextRequest) (*ai.GenerateTextResponse, error) {
	messages := make([]sdk.ChatCompletionMessage, 0, 2)
	
	if req.System != "" {
		messages = append(messages, sdk.ChatCompletionMessage{
			Role:    sdk.ChatMessageRoleSystem,
			Content: req.System,
		})
	}
	messages = append(messages, sdk.ChatCompletionMessage{
		Role:    sdk.ChatMessageRoleUser,
		Content: req.Prompt,
	})

	model := req.Model
	if model == "" {
		model = sdk.GPT3Dot5Turbo
	}

	resp, err := p.client.CreateChatCompletion(ctx, sdk.ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	})
	
	if err != nil {
		return nil, err
	}

	text := ""
	if len(resp.Choices) > 0 {
		text = resp.Choices[0].Message.Content
	}

	return &ai.GenerateTextResponse{
		Text:         text,
		Provider:     ai.ProviderOpenAI,
		ModelUsed:    resp.Model,
		TokensPrompt: resp.Usage.PromptTokens,
		TokensOutput: resp.Usage.CompletionTokens,
	}, nil
}
