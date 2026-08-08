package provider

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/cloudstorex/backend/internal/storage"
)

// RetryConfig holds the configuration for exponential backoff retries.
type RetryConfig struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

// DefaultRetryConfig returns a sensible default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries: 3,
		BaseDelay:  100 * time.Millisecond,
		MaxDelay:   2 * time.Second,
	}
}

// IsRetryable determines if a domain error is temporary and can be retried.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	// Extract underlying DomainError if wrapped
	var domainErr *storage.DomainError
	if errors.As(err, &domainErr) {
		err = domainErr.Err
	}

	// Retry only specific errors
	if errors.Is(err, storage.ErrProviderUnavailable) {
		return true
	}
	
	// Expand if needed
	return false
}

// WithRetry executes the given operation with exponential backoff and jitter.
func WithRetry(ctx context.Context, config RetryConfig, op func() error) error {
	var err error
	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		err = op()
		if err == nil {
			return nil
		}

		if !IsRetryable(err) || attempt == config.MaxRetries {
			return err
		}

		// Calculate backoff with jitter
		delay := float64(config.BaseDelay) * float64(int(1)<<attempt)
		if delay > float64(config.MaxDelay) {
			delay = float64(config.MaxDelay)
		}
		
		// Add up to 20% jitter
		jitter := (rand.Float64() * 0.2) * delay
		finalDelay := time.Duration(delay + jitter)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(finalDelay):
			// continue retry
		}
	}
	return err
}
