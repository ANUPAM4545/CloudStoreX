import { create } from "zustand";
import { UploadItem, UploadStatus } from "./types";

interface UploadStoreState {
  queue: UploadItem[];
  addUpload: (file: File, bucket: string, key: string, controller: AbortController) => UploadItem;
  updateProgress: (id: string, progress: number, status?: UploadStatus, error?: string) => void;
  cancelUpload: (id: string) => void;
  clearCompleted: () => void;
}

export const useUploadStore = create<UploadStoreState>((set, get) => ({
  queue: [],

  addUpload: (file, bucket, key, controller) => {
    const item: UploadItem = {
      id: "upl_" + Math.random().toString(36).substring(2, 11),
      file,
      bucket,
      key,
      progress: 0,
      status: "uploading",
      controller,
    };
    set((state) => ({ queue: [item, ...state.queue] }));
    return item;
  },

  updateProgress: (id, progress, status, error) => {
    set((state) => ({
      queue: state.queue.map((item) => {
        if (item.id !== id) return item;
        return {
          ...item,
          progress,
          status: status || item.status,
          error: error || item.error,
        };
      }),
    }));
  },

  cancelUpload: (id) => {
    const item = get().queue.find((i) => i.id === id);
    if (item && item.controller) {
      item.controller.abort();
    }
    set((state) => ({
      queue: state.queue.map((i) =>
        i.id === id ? { ...i, status: "cancelled", progress: 0 } : i
      ),
    }));
  },

  clearCompleted: () => {
    set((state) => ({
      queue: state.queue.filter(
        (item) => item.status === "uploading" || item.status === "pending"
      ),
    }));
  },
}));
