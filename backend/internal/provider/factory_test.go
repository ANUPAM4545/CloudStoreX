package provider

import (
	"context"
	"testing"

	"github.com/cloudstorex/backend/internal/storage"
)

func TestFactory_RegisterCreatorAndCreate(t *testing.T) {
	deps := &FactoryDeps{}
	factory := NewFactory(deps)

	mockCreator := func(ctx context.Context, info *storage.ProviderInfo, d *FactoryDeps) (storage.StorageProvider, error) {
		if d != deps {
			t.Errorf("expected dependencies to be passed to creator")
		}
		return &MockProvider{}, nil
	}

	err := factory.RegisterCreator(storage.ProviderTypeDefault, mockCreator)
	if err != nil {
		t.Fatalf("expected nil error on RegisterCreator, got %v", err)
	}

	info := &storage.ProviderInfo{
		ID:   "test-id",
		Name: "Test Provider",
		Type: storage.ProviderTypeDefault,
	}

	prov, err := factory.Create(context.Background(), info)
	if err != nil {
		t.Fatalf("expected nil error on Create, got %v", err)
	}
	if prov == nil {
		t.Fatalf("expected valid provider instance, got nil")
	}
}

func TestFactory_UnknownProviderType(t *testing.T) {
	factory := NewFactory(nil)
	info := &storage.ProviderInfo{
		ID:   "test-id",
		Type: "unknown-type",
	}

	_, err := factory.Create(context.Background(), info)
	if err == nil {
		t.Fatalf("expected error on unknown provider type, got nil")
	}
}
