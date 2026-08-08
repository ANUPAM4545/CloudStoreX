package engine

import (
	"testing"

	"github.com/cloudstorex/backend/internal/policy/model"
)

func TestEvaluateNode_LeafOperators(t *testing.T) {
	ctx := &model.EvaluationContext{
		Size:     1024,
		MimeType: "image/png",
		Region:   "us-east-1",
		Tags: map[string]string{
			"env": "production",
		},
	}

	tests := []struct {
		name    string
		leaf    model.ConditionLeaf
		want    bool
		wantErr bool
	}{
		{
			name: "eq operator match",
			leaf: model.ConditionLeaf{Field: "region", Operator: model.OpEq, Value: "us-east-1"},
			want: true,
		},
		{
			name: "eq operator mismatch",
			leaf: model.ConditionLeaf{Field: "region", Operator: model.OpEq, Value: "us-west-1"},
			want: false,
		},
		{
			name: "gt operator match",
			leaf: model.ConditionLeaf{Field: "object.size", Operator: model.OpGt, Value: float64(1000)},
			want: true,
		},
		{
			name: "regex match",
			leaf: model.ConditionLeaf{Field: "object.mime_type", Operator: model.OpRegex, Value: "^image/.*"},
			want: true,
		},
		{
			name: "tags match",
			leaf: model.ConditionLeaf{Field: "tags.env", Operator: model.OpEq, Value: "production"},
			want: true,
		},
		{
			name: "missing tag fallback",
			leaf: model.ConditionLeaf{Field: "tags.nonexistent", Operator: model.OpEq, Value: "foo"},
			want: false, // Doesn't match if it doesn't exist
		},
		{
			name: "exists true",
			leaf: model.ConditionLeaf{Field: "tags.env", Operator: model.OpExists, Value: true},
			want: true,
		},
		{
			name: "exists false",
			leaf: model.ConditionLeaf{Field: "tags.nonexistent", Operator: model.OpExists, Value: false},
			want: true,
		},
		{
			name: "in operator",
			leaf: model.ConditionLeaf{Field: "region", Operator: model.OpIn, Value: []interface{}{"us-west-1", "us-east-1"}},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluateLeaf(ctx, &tt.leaf)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateLeaf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("evaluateLeaf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluateNode_AST(t *testing.T) {
	ctx := &model.EvaluationContext{
		Size:   5000,
		Region: "us-east-1",
	}

	// AND ( size > 1000, region == "us-east-1" )
	node := &model.ConditionNode{
		Type: model.NodeAnd,
		Conditions: []model.ConditionNode{
			{
				Type: model.NodeLeaf,
				Leaf: &model.ConditionLeaf{
					Field:    "object.size",
					Operator: model.OpGt,
					Value:    float64(1000),
				},
			},
			{
				Type: model.NodeLeaf,
				Leaf: &model.ConditionLeaf{
					Field:    "region",
					Operator: model.OpEq,
					Value:    "us-east-1",
				},
			},
		},
	}

	matched, err := evaluateNode(ctx, node)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !matched {
		t.Fatalf("expected node to match")
	}
}
