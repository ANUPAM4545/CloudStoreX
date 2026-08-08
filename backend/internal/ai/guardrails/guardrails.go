package guardrails

import (
	"context"
	"errors"
	"strings"
)

type GuardrailService interface {
	ValidateInput(ctx context.Context, input string) error
	ValidateOutput(ctx context.Context, output string) error
}

type guardrailService struct{}

func NewGuardrailService() GuardrailService {
	return &guardrailService{}
}

func (s *guardrailService) ValidateInput(ctx context.Context, input string) error {
	// E.g. Check for prompt injection, PII masking, secret filtering
	if strings.Contains(strings.ToLower(input), "ignore previous instructions") {
		return errors.New("potential prompt injection detected")
	}
	return nil
}

func (s *guardrailService) ValidateOutput(ctx context.Context, output string) error {
	// E.g. Hallucination confidence scoring, sanitization
	if strings.Contains(strings.ToLower(output), "password") {
		// Just a mock check
		return errors.New("sensitive information detected in output")
	}
	return nil
}
