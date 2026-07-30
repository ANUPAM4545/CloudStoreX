"use client";

import React, { useState } from "react";
import { useUploadStore } from "./store";
import { formatBytes } from "@/features/objects/mappers";
import {
  CheckCircle2,
  AlertCircle,
  X,
  ChevronUp,
  ChevronDown,
  UploadCloud,
  Ban,
} from "lucide-react";


export function UploadProgress() {
  const queue = useUploadStore((state) => state.queue);
  const cancelUpload = useUploadStore((state) => state.cancelUpload);
  const clearCompleted = useUploadStore((state) => state.clearCompleted);
  const [isMinimized, setIsMinimized] = useState(false);

  if (queue.length === 0) return null;

  const activeUploads = queue.filter(
    (item) => item.status === "uploading" || item.status === "pending"
  );
  const completedCount = queue.filter((item) => item.status === "completed").length;


  return (
    <div className="fixed bottom-4 right-4 z-50 w-80 sm:w-96 rounded-xl bg-card border shadow-xl overflow-hidden animate-slide-up">
      {/* Header */}
      <div
        onClick={() => setIsMinimized(!isMinimized)}
        className="flex items-center justify-between bg-primary px-4 py-3 text-primary-foreground cursor-pointer select-none"
      >
        <div className="flex items-center gap-2">
          <UploadCloud size={18} className="animate-pulse" />
          <span className="text-sm font-semibold">
            {activeUploads.length > 0
              ? `Uploading ${activeUploads.length} ${
                  activeUploads.length === 1 ? "file" : "files"
                }...`
              : completedCount > 0
              ? `${completedCount} upload${completedCount === 1 ? "" : "s"} completed`
              : "Upload Manager"}
          </span>
        </div>
        <div className="flex items-center gap-1">
          {queue.some((i) => i.status !== "uploading" && i.status !== "pending") && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                clearCompleted();
              }}
              className="text-xs text-primary-foreground/80 hover:text-white underline mr-1"
            >
              Clear
            </button>
          )}
          <button className="text-primary-foreground/80 hover:text-white">
            {isMinimized ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
          </button>
        </div>
      </div>

      {/* Upload Items List */}
      {!isMinimized && (
        <div className="max-h-64 overflow-y-auto divide-y bg-card">
          {queue.map((item) => (
            <div key={item.id} className="p-3 text-xs space-y-1.5">
              <div className="flex items-center justify-between gap-2">
                <span className="font-medium text-foreground truncate max-w-[200px]" title={item.key}>
                  {item.key}
                </span>
                <div className="flex items-center gap-1.5 shrink-0">
                  {item.status === "uploading" && (
                    <span className="text-muted-foreground">{item.progress}%</span>
                  )}
                  {item.status === "completed" && (
                    <CheckCircle2 size={15} className="text-emerald-500" />
                  )}
                  {item.status === "error" && (
                    <span title={item.error}>
                      <AlertCircle size={15} className="text-rose-500" />
                    </span>
                  )}
                  {item.status === "cancelled" && (
                    <span title="Cancelled">
                      <Ban size={15} className="text-amber-500" />
                    </span>
                  )}
                  {item.status === "uploading" && (
                    <button
                      onClick={() => cancelUpload(item.id)}
                      className="text-muted-foreground hover:text-destructive"
                      title="Cancel upload"
                    >
                      <X size={14} />
                    </button>
                  )}
                </div>
              </div>

              {item.status === "uploading" && (
                <div className="w-full h-1.5 bg-muted rounded-full overflow-hidden">
                  <div
                    className="h-full bg-primary transition-all duration-300"
                    style={{ width: `${item.progress}%` }}
                  />
                </div>
              )}

              <div className="flex items-center justify-between text-[11px] text-muted-foreground">
                <span>{formatBytes(item.file.size)}</span>
                {item.error && (
                  <span className="text-rose-500 truncate max-w-[180px]">{item.error}</span>
                )}
                {item.status === "completed" && (
                  <span className="text-emerald-600 dark:text-emerald-400 font-medium">Uploaded</span>
                )}
                {item.status === "cancelled" && <span>Cancelled</span>}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
