package jobs

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudstorex/backend/internal/jobs/model"
)

// Handler processes a specific type of background job.
type Handler interface {
	Handle(ctx context.Context, job *model.Job) error
}

// HandlerFunc allows using a regular function as a Handler.
type HandlerFunc func(ctx context.Context, job *model.Job) error

func (f HandlerFunc) Handle(ctx context.Context, job *model.Job) error {
	return f(ctx, job)
}

// Registry manages the mapping of job types to their respective handlers.
type Registry interface {
	Register(jobType string, handler Handler)
	Get(jobType string) (Handler, error)
}

type defaultRegistry struct {
	mu       sync.RWMutex
	handlers map[string]Handler
}

// NewRegistry creates a new concurrency-safe job handler registry.
func NewRegistry() Registry {
	return &defaultRegistry{
		handlers: make(map[string]Handler),
	}
}

func (r *defaultRegistry) Register(jobType string, handler Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[jobType] = handler
}

func (r *defaultRegistry) Get(jobType string) (Handler, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	handler, ok := r.handlers[jobType]
	if !ok {
		return nil, fmt.Errorf("no handler registered for job type: %s", jobType)
	}
	return handler, nil
}
