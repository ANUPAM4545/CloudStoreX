export type UploadStatus = "pending" | "uploading" | "completed" | "error" | "cancelled";

export interface UploadItem {
  id: string;
  file: File;
  bucket: string;
  key: string;
  progress: number;
  status: UploadStatus;
  error?: string;
  controller?: AbortController;
}
