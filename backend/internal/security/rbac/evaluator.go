package rbac

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

type Evaluator interface {
	HasPermission(ctx context.Context, subjectID uuid.UUID, action string, scope string, scopeID uuid.UUID) (bool, error)
}

type evaluator struct {
	// Repo to fetch bindings and roles
}

func NewEvaluator() Evaluator {
	return &evaluator{}
}

func (e *evaluator) HasPermission(ctx context.Context, subjectID uuid.UUID, action string, scope string, scopeID uuid.UUID) (bool, error) {
	// In reality: 
	// 1. Fetch RoleBindings for subjectID matching the scope/scopeID (or inherited parent scopes)
	// 2. Fetch the corresponding Roles
	// 3. Evaluate if any role contains the 'action' (supporting wildcards like storage:object:*)
	
	// Mock implementation for the engine test
	if action == "" {
		return false, nil
	}
	
	return matchPermission("storage:object:*", action), nil
}

func matchPermission(granted, requested string) bool {
	if granted == "*" {
		return true
	}
	if strings.HasSuffix(granted, "*") {
		prefix := strings.TrimSuffix(granted, "*")
		return strings.HasPrefix(requested, prefix)
	}
	return granted == requested
}
