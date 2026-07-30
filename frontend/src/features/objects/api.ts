import { apiClient } from "@/lib/api/client";
import { APIResponse } from "@/lib/api/types";
import { ObjectDTO } from "./types";

export const objectsApi = {
  async listObjects(bucket: string, prefix?: string): Promise<ObjectDTO[]> {
    const query = prefix ? `?prefix=${encodeURIComponent(prefix)}` : "";
    const res = await apiClient.get<any, APIResponse<ObjectDTO[]>>(
      `/storage/buckets/${encodeURIComponent(bucket)}/objects${query}`
    );
    if (Array.isArray(res)) {
      return res as unknown as ObjectDTO[];
    }
    if (res.data && Array.isArray(res.data)) {
      return res.data;
    }
    if ((res as any).objects && Array.isArray((res as any).objects)) {
      return (res as any).objects;
    }
    return [];
  },

  async deleteObject(bucket: string, key: string): Promise<void> {
    const encodedKey = key
      .split("/")
      .map((part) => encodeURIComponent(part))
      .join("/");
    await apiClient.delete(
      `/storage/buckets/${encodeURIComponent(bucket)}/objects/${encodedKey}`
    );
  },

  async objectExists(bucket: string, key: string): Promise<boolean> {
    const encodedKey = key
      .split("/")
      .map((part) => encodeURIComponent(part))
      .join("/");
    try {
      await apiClient.head(
        `/storage/buckets/${encodeURIComponent(bucket)}/objects/${encodedKey}`
      );
      return true;
    } catch {
      return false;
    }
  },

  async downloadObject(bucket: string, key: string): Promise<void> {
    const encodedKey = key
      .split("/")
      .map((part) => encodeURIComponent(part))
      .join("/");
    const res = await apiClient.get(
      `/storage/buckets/${encodeURIComponent(bucket)}/objects/${encodedKey}`,
      { responseType: "blob" }
    );

    // Create download link and click it
    const blob = new Blob([(res as any) instanceof Blob ? res : (res as any).data || res]);
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = key.split("/").pop() || key;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
  },
};
