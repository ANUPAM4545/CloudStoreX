package provider

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/cloudstorex/backend/internal/storage"
)

var (
	ErrUnknownProviderType = errors.New("unknown storage provider type")
	ErrCreatorExists       = errors.New("provider creator already registered for type")
)

// FactoryDeps holds dependency-injected shared services required by provider implementations.
// Future providers will receive these dependencies automatically without changing the Factory contract.
type FactoryDeps struct {
	Logger     *slog.Logger
	HTTPClient *http.Client
	// Additional shared services (e.g., CredentialVault, CacheClient, MetricsCollector) can be added here.
}

// ProviderCreator is a constructor function signature for instantiating a concrete storage provider.
type ProviderCreator func(ctx context.Context, info *storage.ProviderInfo, deps *FactoryDeps) (storage.StorageProvider, error)

// Factory creates StorageProvider instances from provider configuration metadata.
type Factory interface {
	RegisterCreator(providerType storage.ProviderType, creator ProviderCreator) error
	Create(ctx context.Context, info *storage.ProviderInfo) (storage.StorageProvider, error)
}

type defaultFactory struct {
	mu       sync.RWMutex
	creators map[storage.ProviderType]ProviderCreator
	deps     *FactoryDeps
}

// NewFactory creates a new Provider Factory with dependency injection.
func NewFactory(deps *FactoryDeps) Factory {
	if deps == nil {
		deps = &FactoryDeps{
			Logger:     slog.Default(),
			HTTPClient: http.DefaultClient,
		}
	}
	return &defaultFactory{
		creators: make(map[storage.ProviderType]ProviderCreator),
		deps:     deps,
	}
}

func (f *defaultFactory) RegisterCreator(providerType storage.ProviderType, creator ProviderCreator) error {
	if creator == nil {
		return errors.New("provider creator cannot be nil")
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if _, exists := f.creators[providerType]; exists {
		return fmt.Errorf("%w: %s", ErrCreatorExists, providerType)
	}

	f.creators[providerType] = creator
	return nil
}

func (f *defaultFactory) Create(ctx context.Context, info *storage.ProviderInfo) (storage.StorageProvider, error) {
	if info == nil {
		return nil, errors.New("provider info cannot be nil")
	}

	f.mu.RLock()
	creator, exists := f.creators[info.Type]
	f.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrUnknownProviderType, info.Type)
	}

	return creator(ctx, info, f.deps)
}
