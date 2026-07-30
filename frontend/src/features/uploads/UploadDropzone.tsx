"use client";

import React, { useRef, useState } from "react";
import { toast } from "sonner";
import { UploadCloud, FileUp } from "lucide-react";
import { Button } from "@/components/ui/button";
import { uploadsApi, MAX_UPLOAD_SIZE_BYTES, MAX_UPLOAD_SIZE_MB } from "./api";
import { useUploadStore } from "./store";
import { useQueryClient } from "@tanstack/react-query";

export interface UploadDropzoneProps {
  bucket: string;
  prefix?: string;
  onUploadSuccess?: () => void;
  compact?: boolean;
}

export function UploadDropzone({
  bucket,
  prefix = "",
  onUploadSuccess,
  compact = false,
}: UploadDropzoneProps) {
  const [isDragging, setIsDragging] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const addUpload = useUploadStore((state) => state.addUpload);
  const updateProgress = useUploadStore((state) => state.updateProgress);
  const queryClient = useQueryClient();

  const handleFiles = async (files: FileList | File[]) => {
    const fileArray = Array.from(files);
    if (fileArray.length === 0) return;

    for (const file of fileArray) {
      if (file.size > MAX_UPLOAD_SIZE_BYTES) {
        toast.error(`File "${file.name}" exceeds 100 MB limit`, {
          description: `Maximum allowed size is ${MAX_UPLOAD_SIZE_MB} MB.`,
        });
        continue;
      }

      const key = prefix + file.name;
      const controller = new AbortController();
      const uploadItem = addUpload(file, bucket, key, controller);

      try {
        await uploadsApi.uploadFile(
          bucket,
          file,
          key,
          (percent) => {
            updateProgress(uploadItem.id, percent, "uploading");
          },
          controller
        );

        updateProgress(uploadItem.id, 100, "completed");
        toast.success("Uploaded " + file.name);
        queryClient.invalidateQueries({ queryKey: ["objects", bucket] });
        if (onUploadSuccess) onUploadSuccess();
      } catch (err: any) {
        if (err?.name === "AbortError" || uploadItem.status === "cancelled") {
          updateProgress(uploadItem.id, 0, "cancelled");
        } else {
          const errorMsg = err?.message || "Failed to upload file";
          updateProgress(uploadItem.id, 0, "error", errorMsg);
          toast.error("Upload failed for " + file.name, { description: errorMsg });
        }
      }
    }

    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  const onDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  };

  const onDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  };

  const onDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      handleFiles(e.dataTransfer.files);
    }
  };

  if (compact) {
    return (
      <>
        <input
          ref={fileInputRef}
          type="file"
          multiple
          onChange={(e) => e.target.files && handleFiles(e.target.files)}
          className="hidden"
        />
        <Button
          onClick={() => fileInputRef.current?.click()}
          className="h-9 gap-1.5 font-semibold"
        >
          <FileUp size={16} />
          <span>Upload Files</span>
        </Button>
      </>
    );
  }

  return (
    <div
      onDragOver={onDragOver}
      onDragLeave={onDragLeave}
      onDrop={onDrop}
      onClick={() => fileInputRef.current?.click()}
      className={`relative cursor-pointer rounded-xl border-2 border-dashed p-6 transition-all text-center flex flex-col items-center justify-center gap-2 ${
        isDragging
          ? "border-primary bg-primary/5 text-primary scale-[1.01]"
          : "border-muted-foreground/25 hover:border-primary/50 bg-card hover:bg-muted/20"
      }`}
    >
      <input
        ref={fileInputRef}
        type="file"
        multiple
        onChange={(e) => e.target.files && handleFiles(e.target.files)}
        className="hidden"
      />
      <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-primary/10 text-primary">
        <UploadCloud size={24} />
      </div>
      <div>
        <p className="text-sm font-semibold text-foreground">
          {isDragging ? "Drop files to upload" : "Click or drag files to upload"}
        </p>
        <p className="text-xs text-muted-foreground mt-0.5">
          Max {MAX_UPLOAD_SIZE_MB} MB per file • Streaming multipart storage pipeline
        </p>
      </div>
    </div>
  );
}
