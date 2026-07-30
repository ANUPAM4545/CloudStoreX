package provider

import (
	"errors"
	"fmt"
	"sync"

	"github.com/cloudstorex/backend/internal/storage"
)

var (
	ErrProviderExists   = errors.New("provider already registered in registry")
	ErrProviderNotFound = storage.ErrProviderNotFound
)

// Registry manages registered StorageProvider instances in a thread-safe manner.
type Registry interface {
	Register(id string, provider storage.StorageProvider) error
	Unregister(id string)
	Get(id string) (storage.StorageProvider, error)
	List() []string
}

type memoryRegistry struct {
	mu        sync.RWMutex
	providers map[string]storage.StorageProvider
}

// NewRegistry creates a new thread-safe Provider Registry.
func NewRegistry() Registry {
	return &memoryRegistry{
		providers: make(map[string]storage.StorageProvider),
	}
}

func (r *memoryRegistry) Register(id string, provider storage.StorageProvider) error {
	if id == "" {
		return errors.New("provider id cannot be empty")
	}
	if provider == nil {
		return errors.New("provider instance cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.providers[id]; exists {
		return fmt.Errorf("%w: %s", ErrProviderExists, id)
	}

	r.providers[id] = provider
	return nil
}

func (r *memoryRegistry) Unregister(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.providers, id)
}

func (r *memoryRegistry) Get(id string) (storage.StorageProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, exists := r.providers[id]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrProviderNotFound, id)
	}

	return provider, nil
}

func (r *memoryRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0, len(r.providers))
	for id := range r.providers {
		ids = append(ids, id)
	}
	return ids
}
