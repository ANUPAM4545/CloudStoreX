package zerotrust

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type RequestContext struct {
	IPAddress string
	UserAgent string
	DeviceID  string
	GeoCity   string
	GeoCountry string
}

type Evaluator interface {
	EvaluateTrust(ctx context.Context, userID uuid.UUID, reqContext RequestContext) (bool, error)
}

type evaluator struct{}

func NewEvaluator() Evaluator {
	return &evaluator{}
}

func (e *evaluator) EvaluateTrust(ctx context.Context, userID uuid.UUID, reqContext RequestContext) (bool, error) {
	// Zero Trust evaluation involves:
	// 1. Device Trust (is it a managed device?)
	// 2. IP Restrictions (is it on a trusted network?)
	// 3. Geo Restrictions (impossible travel / blocked countries)
	// 4. Risk Score calculation (e.g. login failures recently)
	
	// Mock logic: Block requests from untrusted countries
	if reqContext.GeoCountry == "XX" {
		return false, errors.New("zero trust policy denied: blocked geography")
	}

	return true, nil
}
