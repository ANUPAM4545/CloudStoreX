export type ProviderStatus = 'PENDING' | 'VALIDATING' | 'READY' | 'DISABLED' | 'FAILED';
export type ProviderHealth = 'HEALTHY' | 'DEGRADED' | 'UNAVAILABLE' | 'UNKNOWN';

export interface ProviderCapabilities {
  multipart_upload: boolean;
  versioning: boolean;
  object_tagging: boolean;
  pre_signed_urls: boolean;
  encryption: boolean;
  object_lock: boolean;
  lifecycle_rules: boolean;
  replication: boolean;
}

export interface Provider {
  id: string;
  workspace_id: string;
  provider_name: string;
  provider_type: string;
  endpoint?: string;
  region?: string;
  bucket_prefix?: string;
  is_default: boolean;
  is_enabled: boolean;
  status: ProviderStatus;
  health: ProviderHealth;
  last_health_check?: string;
  latency_ms: number;
  last_error?: string;
  capabilities: ProviderCapabilities;
  created_at: string;
  updated_at: string;
}

export interface CreateProviderRequest {
  workspace_id: string;
  provider_name: string;
  provider_type: string;
  endpoint?: string;
  region?: string;
  bucket_prefix?: string;
  is_default: boolean;
  credential_ref?: string;
  capabilities: ProviderCapabilities;
}

export interface UpdateProviderRequest {
  endpoint?: string;
  region?: string;
  bucket_prefix?: string;
  credential_ref?: string;
  capabilities?: ProviderCapabilities;
}
