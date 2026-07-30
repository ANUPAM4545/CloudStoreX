import { apiClient } from "@/lib/api/client";
import { APIResponse } from "@/lib/api/types";
import { BucketDTO } from "./types";

export const bucketsApi = {
  async listBuckets(): Promise<BucketDTO[]> {
    const res = await apiClient.get<any, APIResponse<BucketDTO[]>>("/storage/buckets");
    if (Array.isArray(res)) {
      return res as unknown as BucketDTO[];
    }
    if (res.data && Array.isArray(res.data)) {
      return res.data;
    }
    if ((res as any).buckets && Array.isArray((res as any).buckets)) {
      return (res as any).buckets;
    }
    return [];
  },

  async createBucket(name: string): Promise<BucketDTO> {
    const res = await apiClient.post<any, APIResponse<BucketDTO>>("/storage/buckets", {
      bucket: name,
    });
    return res.data || { name, created_at: new Date().toISOString() };
  },

  async deleteBucket(name: string): Promise<void> {
    await apiClient.delete(`/storage/buckets/${encodeURIComponent(name)}`);
  },
};
