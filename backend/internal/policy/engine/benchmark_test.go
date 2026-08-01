package engine

import (
	"context"
	"testing"

	"github.com/cloudstorex/backend/internal/policy/model"
	"github.com/cloudstorex/backend/internal/policy/rules"
	"github.com/google/uuid"
)

type benchMockRepo struct {
	policies []*model.Policy
}

func (m *benchMockRepo) ListEnabled(ctx context.Context, workspaceID string) ([]*model.Policy, error) {
	return m.policies, nil
}

type benchMockResolver struct{}

func (m *benchMockResolver) GetDefaultProviderID(ctx context.Context, workspaceID string) (string, error) {
	return "aws-s3-default", nil
}

type benchMockPublisher struct{}

func (m *benchMockPublisher) PublishRoutingDecision(ctx context.Context, decision *model.RoutingDecision) error {
	return nil
}

// BenchmarkPolicyEngine_Resolve benchmarks rule evaluation and provider routing resolution.
func BenchmarkPolicyEngine_Resolve(b *testing.B) {
	reg := rules.NewRegistry()
	_ = reg.Register(&rules.ObjectSizeRule{})
	eval := NewEvaluator(reg)

	policies := []*model.Policy{
		{
			RuleType:   model.RuleTypeObjectSize,
			Priority:   10,
			Conditions: []byte(`{"min_size": 0, "max_size": 10000000}`),
			Actions:    []byte(`{"provider_id": "minio"}`),
			Enabled:    true,
		},
	}

	repo := &benchMockRepo{policies: policies}
	eng := NewPolicyEngine(repo, eval, &benchMockResolver{}, &benchMockPublisher{})

	ctx := context.Background()
	evalCtx := &model.EvaluationContext{
		WorkspaceID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Size:        2048,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = eng.Resolve(ctx, evalCtx, "upload")
	}
}

// BenchmarkPolicyEngine_ResolveParallel benchmarks thread-safety and concurrent throughput.
func BenchmarkPolicyEngine_ResolveParallel(b *testing.B) {
	reg := rules.NewRegistry()
	_ = reg.Register(&rules.ObjectSizeRule{})
	eval := NewEvaluator(reg)

	policies := []*model.Policy{
		{
			RuleType:   model.RuleTypeObjectSize,
			Priority:   10,
			Conditions: []byte(`{"min_size": 0, "max_size": 10000000}`),
			Actions:    []byte(`{"provider_id": "minio"}`),
			Enabled:    true,
		},
	}

	repo := &benchMockRepo{policies: policies}
	eng := NewPolicyEngine(repo, eval, &benchMockResolver{}, &benchMockPublisher{})

	ctx := context.Background()
	evalCtx := &model.EvaluationContext{
		WorkspaceID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Size:        2048,
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, _ = eng.Resolve(ctx, evalCtx, "upload")
		}
	})
}
