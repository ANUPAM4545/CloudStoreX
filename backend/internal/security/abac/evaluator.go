package abac

import (
	"context"
	"fmt"
)

type ResourceAttributes struct {
	Type     string
	Provider string
	Region   string
	Bucket   string
	Tags     map[string]string
}

type ContextAttributes struct {
	IPAddress string
	Time      string // e.g. UTC hour
	DeviceID  string
}

type Evaluator interface {
	Evaluate(ctx context.Context, resource ResourceAttributes, environment ContextAttributes) (bool, error)
}

type evaluator struct{}

func NewEvaluator() Evaluator {
	return &evaluator{}
}

func (e *evaluator) Evaluate(ctx context.Context, resource ResourceAttributes, environment ContextAttributes) (bool, error) {
	// In reality, this would evaluate a JSON-based policy containing conditions
	// e.g. "Allow if resource.Provider == 'aws' AND environment.IPAddress in '192.168.0.0/16'"
	
	// Mock logic: block access to a specific mock region to test evaluation
	if resource.Region == "restricted-region" {
		return false, fmt.Errorf("ABAC policy denied: restricted region")
	}

	return true, nil
}
