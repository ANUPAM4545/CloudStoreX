package provider

import (
	"errors"
)

type Manager interface {
	// ActiveProvider returns the current active provider based on configuration/load.
	ActiveProvider() (AIProvider, error)

	// FallbackProvider returns a fallback provider if the active one fails.
	FallbackProvider() (AIProvider, error)
}

type manager struct {
	registry Registry
	activeID string
}

func NewManager(registry Registry, activeID string) Manager {
	return &manager{
		registry: registry,
		activeID: activeID,
	}
}

func (m *manager) ActiveProvider() (AIProvider, error) {
	if m.activeID == "" {
		return nil, errors.New("no active AI provider configured")
	}
	return m.registry.Get(m.activeID)
}

func (m *manager) FallbackProvider() (AIProvider, error) {
	// For now, return the first available provider that is not the active one
	providers := m.registry.List()
	for _, p := range providers {
		if p.ID() != m.activeID {
			return p, nil
		}
	}
	return nil, errors.New("no fallback AI provider available")
}
