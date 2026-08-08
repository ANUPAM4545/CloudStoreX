package rules

import (
	"encoding/json"
	"testing"

	"github.com/cloudstorex/backend/internal/policy/model"
	"github.com/google/uuid"
)

func TestReplicationRule_Evaluate(t *testing.T) {
	rule := &ReplicationRule{}

	if rule.RuleType() != model.RuleTypeReplication {
		t.Fatalf("expected RuleTypeReplication, got %s", rule.RuleType())
	}

	conds := replicationConditions{
		MinSize: 100,
	}
	condsJSON, _ := json.Marshal(conds)

	actions := replicationActions{
		ReplicaProviders: []string{"aws-s3", "minio"},
		Mode:             "ASYNC",
	}
	actionsJSON, _ := json.Marshal(actions)

	ctx := &model.EvaluationContext{
		WorkspaceID: uuid.New(),
		Size:        50, // smaller than min_size
	}

	res, err := rule.Evaluate(ctx, condsJSON, actionsJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Matched {
		t.Fatalf("expected no match when object smaller than min_size")
	}

	ctx.Size = 500 // matches condition
	res, err = rule.Evaluate(ctx, condsJSON, actionsJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Matched {
		t.Fatalf("expected rule match for valid object size")
	}
	if res.ProviderID != "" {
		t.Fatalf("expected ProviderID to remain empty so primary routing is unmodified, got %s", res.ProviderID)
	}
}
