package provider

import (
	"errors"
)

type Factory interface {
	Create(providerType string, config map[string]string) (AIProvider, error)
}

type factory struct{}

func NewFactory() Factory {
	return &factory{}
}

func (f *factory) Create(providerType string, config map[string]string) (AIProvider, error) {
	switch providerType {
	case "mock":
		// In a real setup, we would return the mock provider from the mock package,
		// but to avoid cyclic imports inside the factory, providers often self-register
		// or the factory accepts builder functions.
		// We will rely on explicit registration for now.
		return nil, errors.New("use explicit registration or builder functions")
	default:
		return nil, errors.New("unsupported ai provider type: " + providerType)
	}
}
