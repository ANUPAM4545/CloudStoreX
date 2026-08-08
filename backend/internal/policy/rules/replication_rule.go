package rules

import (
	"encoding/json"
	"fmt"

	"github.com/cloudstorex/backend/internal/policy/model"
)

// ReplicationRule evaluates replication policies without altering primary storage routing (Epic 13 Phase 2).
type ReplicationRule struct{}

type replicationConditions struct {
	MinSize      int64  `json:"min_size"`
	MaxSize      int64  `json:"max_size"`
	BucketID     string `json:"bucket_id"`
	Region       string `json:"region"`
	StorageClass string `json:"storage_class"`
}

type replicationActions struct {
	ReplicaProviders []string `json:"replica_providers"`
	Mode             string   `json:"mode"` // ASYNC or SYNC
}

func (r *ReplicationRule) RuleType() model.RuleType {
	return model.RuleTypeReplication
}

func (r *ReplicationRule) Evaluate(ctx *model.EvaluationContext, conditionsJSON, actionsJSON []byte) (RuleResult, error) {
	var conds replicationConditions
	if len(conditionsJSON) > 0 {
		if err := json.Unmarshal(conditionsJSON, &conds); err != nil {
			return RuleResult{}, fmt.Errorf("failed to parse replication conditions: %w", err)
		}
	}

	var actions replicationActions
	if len(actionsJSON) > 0 {
		if err := json.Unmarshal(actionsJSON, &actions); err != nil {
			return RuleResult{}, fmt.Errorf("failed to parse replication actions: %w", err)
		}
	}

	if conds.MinSize > 0 && ctx.Size < conds.MinSize {
		return RuleResult{Matched: false, Explanation: "object smaller than min_size"}, nil
	}
	if conds.MaxSize > 0 && ctx.Size > conds.MaxSize {
		return RuleResult{Matched: false, Explanation: "object larger than max_size"}, nil
	}
	if conds.Region != "" && ctx.Region != conds.Region {
		return RuleResult{Matched: false, Explanation: "region mismatch"}, nil
	}

	if len(actions.ReplicaProviders) == 0 {
		return RuleResult{Matched: false, Explanation: "no replica providers specified"}, nil
	}

	explanation := fmt.Sprintf("Matched replication policy targeting replicas %v", actions.ReplicaProviders)
	return RuleResult{
		Matched:     true,
		ProviderID:  "", // does not override primary routing provider
		Explanation: explanation,
	}, nil
}
