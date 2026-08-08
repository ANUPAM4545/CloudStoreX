package provider

import (
	"errors"
	"sync"
)

var (
	ErrProviderNotFound      = errors.New("ai provider not found")
	ErrProviderAlreadyExists = errors.New("ai provider already exists")
)

type Registry interface {
	Register(p AIProvider) error
	Get(id string) (AIProvider, error)
	List() []AIProvider
	Remove(id string)
}

type registry struct {
	mu        sync.RWMutex
	providers map[string]AIProvider
}

func NewRegistry() Registry {
	return &registry{
		providers: make(map[string]AIProvider),
	}
}

func (r *registry) Register(p AIProvider) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.providers[p.ID()]; exists {
		return ErrProviderAlreadyExists
	}
	r.providers[p.ID()] = p
	return nil
}

func (r *registry) Get(id string) (AIProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.providers[id]
	if !exists {
		return nil, ErrProviderNotFound
	}
	return p, nil
}

func (r *registry) List() []AIProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]AIProvider, 0, len(r.providers))
	for _, p := range r.providers {
		list = append(list, p)
	}
	return list
}

func (r *registry) Remove(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.providers, id)
}
