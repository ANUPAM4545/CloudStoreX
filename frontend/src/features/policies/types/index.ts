export type RuleType = "REGION" | "HEALTH" | "CAPABILITY" | "BUCKET" | "OBJECT_SIZE" | "MIME_TYPE" | "DEFAULT";

export interface Policy {
  id: string;
  workspace_id: string;
  name: string;
  description: string;
  priority: number;
  enabled: boolean;
  rule_type: RuleType;
  conditions: any;
  actions: {
    provider_id: string;
  };
  created_at: string;
  updated_at: string;
}

export interface RoutingDecision {
  id: string;
  workspace_id: string;
  object_key: string;
  provider_id: string;
  policy_id?: string;
  matched_rule: string;
  is_fallback: boolean;
  candidate_providers: string[];
  decision_reason: string;
  operation: string;
  latency_ms: number;
  timestamp: string;
}

export interface CreatePolicyRequest {
  workspace_id: string;
  name: string;
  description: string;
  priority: number;
  rule_type: RuleType;
  conditions: any;
  actions: {
    provider_id: string;
  };
}

export interface UpdatePolicyRequest {
  name?: string;
  description?: string;
  priority?: number;
  rule_type?: RuleType;
  conditions?: any;
  actions?: {
    provider_id: string;
  };
}
