package memory

import (
	"context"
	"errors"

	"github.com/cloudstorex/backend/internal/ai/provider"
	"github.com/google/uuid"
)

type MemoryManager interface {
	AddMessage(ctx context.Context, sessionID uuid.UUID, message provider.ChatMessage) error
	GetHistory(ctx context.Context, sessionID uuid.UUID) ([]provider.ChatMessage, error)
	Clear(ctx context.Context, sessionID uuid.UUID) error
}

type memoryManager struct {
	// Simple in-memory store for now. Real implementation uses Redis/PostgreSQL.
	store map[uuid.UUID][]provider.ChatMessage
}

func NewMemoryManager() MemoryManager {
	return &memoryManager{
		store: make(map[uuid.UUID][]provider.ChatMessage),
	}
}

func (m *memoryManager) AddMessage(ctx context.Context, sessionID uuid.UUID, message provider.ChatMessage) error {
	m.store[sessionID] = append(m.store[sessionID], message)
	return nil
}

func (m *memoryManager) GetHistory(ctx context.Context, sessionID uuid.UUID) ([]provider.ChatMessage, error) {
	history, exists := m.store[sessionID]
	if !exists {
		return []provider.ChatMessage{}, errors.New("no history found for session")
	}
	return history, nil
}

func (m *memoryManager) Clear(ctx context.Context, sessionID uuid.UUID) error {
	delete(m.store, sessionID)
	return nil
}
