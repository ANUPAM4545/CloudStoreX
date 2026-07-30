export interface APIResponse<T = any> {
  success: boolean;
  data?: T;
  error?: APIError;
  meta?: {
    page?: number;
    limit?: number;
    total?: number;
  };
}

export interface APIError {
  code: string;
  message: string;
  details?: Record<string, any>;
}

export interface AuthResponse {
  token: string;
  user?: {
    id: string;
    email: string;
    full_name?: string;
  };
}

export interface BucketDTO {
  name: string;
  created_at: string;
}

export interface ObjectMetadataDTO {
  content_type?: string;
  cache_control?: string;
  content_encoding?: string;
  custom_metadata?: Record<string, string>;
}

export interface ObjectDTO {
  key: string;
  bucket: string;
  size: number;
  last_modified: string;
  etag?: string;
  content_type?: string;
  metadata?: ObjectMetadataDTO;
}

export interface UploadResponseDTO {
  bucket: string;
  key: string;
  size: number;
  etag?: string;
}
