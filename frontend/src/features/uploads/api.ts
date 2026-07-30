import { apiClient } from "@/lib/api/client";
import { APIResponse, UploadResponseDTO } from "@/lib/api/types";

export const MAX_UPLOAD_SIZE_MB = 100;
export const MAX_UPLOAD_SIZE_BYTES = MAX_UPLOAD_SIZE_MB * 1024 * 1024;

export const uploadsApi = {
  async uploadFile(
    bucket: string,
    file: File,
    key?: string,
    onProgress?: (percent: number) => void,
    controller?: AbortController
  ): Promise<UploadResponseDTO> {
    if (file.size > MAX_UPLOAD_SIZE_BYTES) {
      throw new Error(
        `File exceeds the maximum allowed limit of ${MAX_UPLOAD_SIZE_MB} MB. (Size: ${(
          file.size /
          1024 /
          1024
        ).toFixed(1)} MB)`
      );
    }

    const fullKey = key || file.name;
    const formData = new FormData();
    formData.append("key", fullKey);
    formData.append("file", file);

    const res = await apiClient.post<any, APIResponse<UploadResponseDTO>>(
      `/storage/buckets/${encodeURIComponent(bucket)}/objects?key=${encodeURIComponent(
        fullKey
      )}`,
      formData,
      {
        headers: {
          "Content-Type": "multipart/form-data",
        },
        signal: controller?.signal,
        onUploadProgress: (progressEvent) => {
          if (progressEvent.total && onProgress) {
            const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total);
            onProgress(percent);
          }
        },
      }
    );

    return res.data || {
      bucket,
      key: fullKey,
      size: file.size,
    };
  },
};
