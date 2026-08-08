package engine

import (
	"context"
	"errors"

	"github.com/cloudstorex/backend/internal/security/abac"
	"github.com/cloudstorex/backend/internal/security/rbac"
	"github.com/cloudstorex/backend/internal/security/zerotrust"
	"github.com/google/uuid"
)

type AuthorizationEngine interface {
	Authorize(ctx context.Context, req AuthorizationRequest) error
}

type AuthorizationRequest struct {
	SubjectID uuid.UUID
	
	// RBAC Context
	Action  string
	Scope   string
	ScopeID uuid.UUID
	
	// ABAC Context
	Resource abac.ResourceAttributes
	
	// Zero Trust Context
	ZeroTrust zerotrust.RequestContext
}

type engine struct {
	rbacEvaluator rbac.Evaluator
	abacEvaluator abac.Evaluator
	ztEvaluator   zerotrust.Evaluator
}

func NewAuthorizationEngine(r rbac.Evaluator, a abac.Evaluator, z zerotrust.Evaluator) AuthorizationEngine {
	return &engine{
		rbacEvaluator: r,
		abacEvaluator: a,
		ztEvaluator:   z,
	}
}

func (e *engine) Authorize(ctx context.Context, req AuthorizationRequest) error {
	// Step 1: Zero Trust Evaluation
	trustOk, err := e.ztEvaluator.EvaluateTrust(ctx, req.SubjectID, req.ZeroTrust)
	if err != nil || !trustOk {
		return errors.New("access denied by Zero Trust policy: " + err.Error())
	}

	// Step 2: RBAC Evaluation
	rbacOk, err := e.rbacEvaluator.HasPermission(ctx, req.SubjectID, req.Action, req.Scope, req.ScopeID)
	if err != nil || !rbacOk {
		return errors.New("access denied by RBAC policy: insufficient permissions")
	}

	// Step 3: ABAC Evaluation
	abacCtx := abac.ContextAttributes{
		IPAddress: req.ZeroTrust.IPAddress,
		DeviceID:  req.ZeroTrust.DeviceID,
	}
	abacOk, err := e.abacEvaluator.Evaluate(ctx, req.Resource, abacCtx)
	if err != nil || !abacOk {
		return errors.New("access denied by ABAC policy: " + err.Error())
	}

	// Authorization Decision: Granted
	return nil
}
