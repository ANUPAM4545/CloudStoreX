"use client";

import React from "react";
import { toast } from "sonner";
import { ConfirmDialog } from "@/components/ui/ConfirmDialog";
import { useDeleteBucket } from "./hooks";
import { BucketViewModel } from "./types";

export interface DeleteBucketConfirmProps {
  bucket: BucketViewModel | null;
  onClose: () => void;
  onDeleted?: () => void;
}

export function DeleteBucketConfirm({ bucket, onClose, onDeleted }: DeleteBucketConfirmProps) {
  const deleteBucketMutation = useDeleteBucket();

  if (!bucket) return null;

  const handleConfirm = async () => {
    try {
      await deleteBucketMutation.mutateAsync(bucket.name);
      toast.success("Bucket deleted", {
        description: `Bucket "${bucket.name}" has been permanently removed.`,
      });
      if (onDeleted) onDeleted();
    } catch (err: any) {
      const errorMsg = err?.message || "Failed to delete bucket. Ensure the bucket is empty.";
      toast.error("Deletion failed", { description: errorMsg });
      throw err; // rethrow so ConfirmDialog stops processing state
    }
  };

  return (
    <ConfirmDialog
      isOpen={!!bucket}
      onClose={onClose}
      onConfirm={handleConfirm}
      title="Delete Storage Bucket"
      description={
        <span>
          Are you sure you want to delete <strong className="text-foreground">{bucket.name}</strong>?
          This action cannot be undone. All files within this bucket must be deleted first.
        </span>
      }
      confirmLabel="Delete Bucket"
      variant="destructive"
      requireConfirmationText={bucket.name}
    />
  );
}
