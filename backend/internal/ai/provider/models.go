package provider

import (
	"time"
)

type GenerateRequest struct {
	Prompt      string
	Model       string
	Temperature float64
	MaxTokens   int
	SystemPrompt string
	Stream      bool
}

type GenerateResponse struct {
	Text         string
	InputTokens  int
	OutputTokens int
	TotalTokens  int
	Cost         float64
	Model        string
	ProviderName string
	Latency      time.Duration
}

type ChatMessageRole string

const (
	RoleSystem    ChatMessageRole = "system"
	RoleUser      ChatMessageRole = "user"
	RoleAssistant ChatMessageRole = "assistant"
)

type ChatMessage struct {
	Role    ChatMessageRole
	Content string
}

type ChatRequest struct {
	Messages    []ChatMessage
	Model       string
	Temperature float64
	MaxTokens   int
	Stream      bool
}

type EmbeddingRequest struct {
	Input string
	Model string
}

type EmbeddingResponse struct {
	Embedding    []float32
	InputTokens  int
	TotalTokens  int
	Cost         float64
	Model        string
	ProviderName string
	Latency      time.Duration
}
