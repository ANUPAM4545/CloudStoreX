package assistant

import (
	"context"

	aictx "github.com/cloudstorex/backend/internal/ai/context"
	"github.com/cloudstorex/backend/internal/ai/memory"
	"github.com/cloudstorex/backend/internal/ai/provider"
	"github.com/google/uuid"
)

type Assistant interface {
	Ask(ctx context.Context, sessionID uuid.UUID, prompt string) (string, error)
}

type assistant struct {
	aiManager provider.Manager
	memory    memory.MemoryManager
	ctxBuild  aictx.Builder
}

func NewAssistant(m provider.Manager, mem memory.MemoryManager, cb aictx.Builder) Assistant {
	return &assistant{
		aiManager: m,
		memory:    mem,
		ctxBuild:  cb,
	}
}

func (a *assistant) Ask(ctx context.Context, sessionID uuid.UUID, prompt string) (string, error) {
	// 1. Get active provider
	activeProvider, err := a.aiManager.ActiveProvider()
	if err != nil {
		return "", err
	}

	// 2. Add user prompt to memory
	userMsg := provider.ChatMessage{
		Role:    provider.RoleUser,
		Content: prompt,
	}
	a.memory.AddMessage(ctx, sessionID, userMsg)

	// 3. Retrieve chat history
	history, err := a.memory.GetHistory(ctx, sessionID)
	if err != nil {
		history = []provider.ChatMessage{userMsg}
	}

	// 4. Optionally inject system context here via a.ctxBuild...
	
	// 5. Send to AI Provider
	req := provider.ChatRequest{
		Messages:    history,
		Model:       "default-assistant-model",
		Temperature: 0.7,
	}

	resp, err := activeProvider.ChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	// 6. Save assistant response to memory
	assistMsg := provider.ChatMessage{
		Role:    provider.RoleAssistant,
		Content: resp.Text,
	}
	a.memory.AddMessage(ctx, sessionID, assistMsg)

	return resp.Text, nil
}
