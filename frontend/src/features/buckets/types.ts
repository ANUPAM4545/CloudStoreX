export interface BucketDTO {
  name: string;
  created_at: string;
}

export interface BucketViewModel {
  name: string;
  createdAtFormatted: string;
  createdAtRaw: string;
}

export interface CreateBucketRequestDTO {
  name: string;
}
