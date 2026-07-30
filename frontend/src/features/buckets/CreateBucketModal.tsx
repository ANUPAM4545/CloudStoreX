"use client";

import React from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { X, FolderPlus } from "lucide-react";
import { createBucketSchema, CreateBucketFormData } from "./schemas";
import { useCreateBucket } from "./hooks";

export interface CreateBucketModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCreated?: (bucketName: string) => void;
}

export function CreateBucketModal({ isOpen, onClose, onCreated }: CreateBucketModalProps) {
  const createBucketMutation = useCreateBucket();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CreateBucketFormData>({
    resolver: zodResolver(createBucketSchema),
    defaultValues: {
      name: "",
    },
  });

  if (!isOpen) return null;

  const onSubmit = async (data: CreateBucketFormData) => {
    try {
      await createBucketMutation.mutateAsync(data.name);
      toast.success("Bucket created successfully", {
        description: `Bucket "${data.name}" is now ready for object storage.`,
      });
      reset();
      onClose();
      if (onCreated) onCreated(data.name);
    } catch (err: any) {
      const errorMsg = err?.message || "Failed to create bucket. A bucket with this name may already exist.";
      toast.error("Failed to create bucket", { description: errorMsg });
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-xs animate-fade-in p-4">
      <div className="relative w-full max-w-md rounded-xl bg-card border shadow-lg p-6 animate-slide-up">
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-muted-foreground hover:text-foreground rounded-md p-1 transition-colors"
        >
          <X size={18} />
        </button>

        <div className="flex items-center gap-3 mb-4">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 text-primary">
            <FolderPlus size={20} />
          </div>
          <div>
            <h3 className="text-lg font-semibold tracking-tight text-foreground">Create Storage Bucket</h3>
            <p className="text-xs text-muted-foreground">Add a new isolated container for your cloud objects.</p>
          </div>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" noValidate>
          <div className="space-y-1.5">
            <label
              htmlFor="bucket-name"
              className="text-xs font-medium uppercase tracking-wider text-muted-foreground"
            >
              Bucket Name
            </label>
            <Input
              id="bucket-name"
              type="text"
              placeholder="e.g. my-production-assets"
              autoFocus
              {...register("name")}
              className={errors.name ? "border-destructive focus-visible:ring-destructive" : ""}
            />
            {errors.name ? (
              <p className="text-xs text-destructive mt-1">{errors.name.message}</p>
            ) : (
              <p className="text-xs text-muted-foreground mt-1">
                Must be 3-63 lowercase alphanumeric characters or hyphens.
              </p>
            )}
          </div>

          <div className="flex items-center justify-end gap-3 mt-6 pt-2 border-t">
            <Button
              type="button"
              variant="outline"
              onClick={onClose}
              disabled={createBucketMutation.isPending}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={createBucketMutation.isPending}
            >
              {createBucketMutation.isPending ? "Creating..." : "Create Bucket"}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
