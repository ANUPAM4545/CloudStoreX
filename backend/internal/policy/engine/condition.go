package engine

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/cloudstorex/backend/internal/policy/model"
)

// evaluateNode recursively evaluates a condition AST node against the evaluation context.
func evaluateNode(ctx *model.EvaluationContext, node *model.ConditionNode) (bool, error) {
	if node == nil {
		return true, nil
	}

	switch node.Type {
	case model.NodeAnd:
		if len(node.Conditions) == 0 {
			return true, nil
		}
		for _, child := range node.Conditions {
			matched, err := evaluateNode(ctx, &child)
			if err != nil {
				return false, err
			}
			if !matched {
				return false, nil
			}
		}
		return true, nil

	case model.NodeOr:
		if len(node.Conditions) == 0 {
			return false, nil
		}
		for _, child := range node.Conditions {
			matched, err := evaluateNode(ctx, &child)
			if err != nil {
				return false, err
			}
			if matched {
				return true, nil
			}
		}
		return false, nil

	case model.NodeLeaf:
		if node.Leaf == nil {
			return true, nil
		}
		return evaluateLeaf(ctx, node.Leaf)

	default:
		return false, fmt.Errorf("unknown condition node type: %s", node.Type)
	}
}

// evaluateLeaf evaluates a single condition leaf against the context.
func evaluateLeaf(ctx *model.EvaluationContext, leaf *model.ConditionLeaf) (bool, error) {
	fieldValue, exists := getFieldValue(ctx, leaf.Field)
	
	// Handle 'exists' operator specifically
	if leaf.Operator == model.OpExists {
		wantExists, ok := leaf.Value.(bool)
		if !ok {
			return false, fmt.Errorf("value for 'exists' operator must be boolean")
		}
		return exists == wantExists, nil
	}

	// If field doesn't exist and we're not checking for existence, fail gracefully
	if !exists {
		return false, nil
	}

	return evaluateOperator(fieldValue, leaf.Operator, leaf.Value)
}

// getFieldValue extracts a value from the EvaluationContext based on field path (e.g., "object.size", "tags.env").
func getFieldValue(ctx *model.EvaluationContext, field string) (interface{}, bool) {
	switch {
	case field == "workspace_id":
		return ctx.WorkspaceID.String(), true
	case field == "organization_id":
		return ctx.OrganizationID.String(), true
	case field == "bucket":
		return ctx.Bucket, true
	case field == "object_key":
		return ctx.ObjectKey, true
	case field == "object.size":
		return float64(ctx.Size), true // Normalize numbers to float64 for generic evaluation
	case field == "object.mime_type":
		return ctx.MimeType, true
	case field == "object.storage_class":
		return ctx.StorageClass, true
	case field == "region":
		return ctx.Region, true
	case strings.HasPrefix(field, "tags."):
		tagKey := strings.TrimPrefix(field, "tags.")
		if ctx.Tags != nil {
			val, ok := ctx.Tags[tagKey]
			return val, ok
		}
		return nil, false
	case strings.HasPrefix(field, "metadata."):
		metaKey := strings.TrimPrefix(field, "metadata.")
		if ctx.Metadata != nil {
			val, ok := ctx.Metadata[metaKey]
			return val, ok
		}
		return nil, false
	default:
		return nil, false
	}
}

// evaluateOperator performs the generic operator comparison.
func evaluateOperator(actual interface{}, op model.ConditionOperator, expected interface{}) (bool, error) {
	switch op {
	case model.OpEq:
		return reflect.DeepEqual(actual, expected), nil
	case model.OpNeq:
		return !reflect.DeepEqual(actual, expected), nil
	case model.OpContains:
		actualStr, ok1 := actual.(string)
		expectedStr, ok2 := expected.(string)
		if ok1 && ok2 {
			return strings.Contains(actualStr, expectedStr), nil
		}
		return false, fmt.Errorf("operator 'contains' requires string operands")
	case model.OpStartsWith:
		actualStr, ok1 := actual.(string)
		expectedStr, ok2 := expected.(string)
		if ok1 && ok2 {
			return strings.HasPrefix(actualStr, expectedStr), nil
		}
		return false, fmt.Errorf("operator 'starts_with' requires string operands")
	case model.OpEndsWith:
		actualStr, ok1 := actual.(string)
		expectedStr, ok2 := expected.(string)
		if ok1 && ok2 {
			return strings.HasSuffix(actualStr, expectedStr), nil
		}
		return false, fmt.Errorf("operator 'ends_with' requires string operands")
	case model.OpRegex:
		actualStr, ok1 := actual.(string)
		expectedStr, ok2 := expected.(string)
		if ok1 && ok2 {
			matched, err := regexp.MatchString(expectedStr, actualStr)
			return matched, err
		}
		return false, fmt.Errorf("operator 'regex' requires string operands")
	case model.OpGt, model.OpGte, model.OpLt, model.OpLte:
		actualF, ok1 := toFloat64(actual)
		expectedF, ok2 := toFloat64(expected)
		if !ok1 || !ok2 {
			return false, fmt.Errorf("numeric operators require numeric operands")
		}
		switch op {
		case model.OpGt:
			return actualF > expectedF, nil
		case model.OpGte:
			return actualF >= expectedF, nil
		case model.OpLt:
			return actualF < expectedF, nil
		case model.OpLte:
			return actualF <= expectedF, nil
		}
	case model.OpIn, model.OpNotIn:
		expectedSlice, ok := expected.([]interface{})
		if !ok {
			return false, fmt.Errorf("operator 'in/not_in' requires an array expected value")
		}
		found := false
		for _, e := range expectedSlice {
			if reflect.DeepEqual(actual, e) {
				found = true
				break
			}
		}
		if op == model.OpIn {
			return found, nil
		}
		return !found, nil
	}
	
	return false, fmt.Errorf("unsupported operator: %s", op)
}

func toFloat64(v interface{}) (float64, bool) {
	switch i := v.(type) {
	case float64:
		return i, true
	case float32:
		return float64(i), true
	case int:
		return float64(i), true
	case int64:
		return float64(i), true
	case int32:
		return float64(i), true
	}
	return 0, false
}
