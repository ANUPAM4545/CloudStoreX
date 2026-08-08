# Authorization

All endpoints route through the `Central Authorization Engine`, which acts as the ultimate gatekeeper evaluating access via:

1. Zero Trust Evaluation
2. RBAC Validation
3. ABAC Constraints

Access is denied if any layer flags the context as non-compliant.
