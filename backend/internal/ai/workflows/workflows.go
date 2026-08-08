package workflows

import (
	"context"
	"errors"
)

type StepContext map[string]interface{}

type Step func(ctx context.Context, input StepContext) (StepContext, error)

type Workflow struct {
	Name  string
	Steps []Step
}

type Engine interface {
	Register(workflow Workflow) error
	Execute(ctx context.Context, name string, initialInput StepContext) (StepContext, error)
}

type engine struct {
	workflows map[string]Workflow
}

func NewEngine() Engine {
	return &engine{
		workflows: make(map[string]Workflow),
	}
}

func (e *engine) Register(w Workflow) error {
	if _, exists := e.workflows[w.Name]; exists {
		return errors.New("workflow already exists")
	}
	e.workflows[w.Name] = w
	return nil
}

func (e *engine) Execute(ctx context.Context, name string, initialInput StepContext) (StepContext, error) {
	w, exists := e.workflows[name]
	if !exists {
		return nil, errors.New("workflow not found")
	}

	currentInput := initialInput
	var err error

	for _, step := range w.Steps {
		currentInput, err = step(ctx, currentInput)
		if err != nil {
			return nil, err
		}
	}

	return currentInput, nil
}
