package tools

import (
	"context"
	"errors"
)

type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, input map[string]interface{}) (interface{}, error)
}

type Registry interface {
	Register(tool Tool) error
	Get(name string) (Tool, error)
	List() []Tool
}

type registry struct {
	tools map[string]Tool
}

func NewRegistry() Registry {
	return &registry{
		tools: make(map[string]Tool),
	}
}

func (r *registry) Register(t Tool) error {
	if _, exists := r.tools[t.Name()]; exists {
		return errors.New("tool already registered")
	}
	r.tools[t.Name()] = t
	return nil
}

func (r *registry) Get(name string) (Tool, error) {
	t, exists := r.tools[name]
	if !exists {
		return nil, errors.New("tool not found")
	}
	return t, nil
}

func (r *registry) List() []Tool {
	var list []Tool
	for _, t := range r.tools {
		list = append(list, t)
	}
	return list
}
