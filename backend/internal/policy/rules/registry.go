package rules

import (
	"fmt"
	"sync"

	"github.com/cloudstorex/backend/internal/policy/model"
)

// Registry manages the registration and lookup of policy rules.
type Registry interface {
	Register(rule Rule) error
	Get(ruleType model.RuleType) (Rule, error)
	List() []Rule
}

type defaultRegistry struct {
	mu    sync.RWMutex
	rules map[model.RuleType]Rule
}

// NewRegistry creates a new rule registry.
func NewRegistry() Registry {
	return &defaultRegistry{
		rules: make(map[model.RuleType]Rule),
	}
}

func (r *defaultRegistry) Register(rule Rule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if rule == nil {
		return fmt.Errorf("cannot register nil rule")
	}

	r.rules[rule.RuleType()] = rule
	return nil
}

func (r *defaultRegistry) Get(ruleType model.RuleType) (Rule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rule, ok := r.rules[ruleType]
	if !ok {
		return nil, fmt.Errorf("rule type %s not found in registry", ruleType)
	}
	return rule, nil
}

func (r *defaultRegistry) List() []Rule {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]Rule, 0, len(r.rules))
	for _, rule := range r.rules {
		list = append(list, rule)
	}
	return list
}
