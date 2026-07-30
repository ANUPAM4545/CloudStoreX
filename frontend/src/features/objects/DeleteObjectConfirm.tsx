"use client";

import React from "react";
import { toast } from "sonner";
import { ConfirmDialog } from "@/components/ui/ConfirmDialog";
import { useDeleteObject } from "./hooks";
import { ObjectViewModel } from "./types";

export interface DeleteObjectConfirmProps {
  bucket: string;
  objectItem: ObjectViewModel | null;
  onClose: () => void;
  onDeleted?: () => void;
}

export function DeleteObjectConfirm({
  bucket,
  objectItem,
  onClose,
  onDeleted,
}: DeleteObjectConfirmProps) {
  const deleteObjectMutation = useDeleteObject(bucket);

  if (!objectItem) return null;

  const handleConfirm = async () => {
    try {
      await deleteObjectMutation.mutateAsync(objectItem.key);
      toast.success("Object deleted", {
        description: `"${objectItem.name}" has been permanently deleted from storage.`,
      });
      if (onDeleted) onDeleted();
    } catch (err: any) {
      const errorMsg = err?.message || "Failed to delete object from storage provider.";
      toast.error("Deletion failed", { description: errorMsg });
      throw err;
    }
  };

  return (
    <ConfirmDialog
      isOpen={!!objectItem}
      onClose={onClose}
      onConfirm={handleConfirm}
      title="Delete Storage Object"
      description={
        <span>
          Are you sure you want to delete <strong className="text-foreground">{objectItem.name}</strong>?
          This action will remove the object from <strong className="text-foreground">{bucket}</strong> and cannot be undone.
        </span>
      }
      confirmLabel="Delete Object"
      variant="destructive"
    />
  );
}
