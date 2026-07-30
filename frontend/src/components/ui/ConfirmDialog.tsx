"use client";

import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { AlertTriangle, X } from "lucide-react";
import { clsx } from "clsx";

export interface ConfirmDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => Promise<void> | void;
  title: string;
  description: React.ReactNode;
  confirmLabel?: string;
  cancelLabel?: string;
  variant?: "destructive" | "default";
  requireConfirmationText?: string;
}

export function ConfirmDialog({
  isOpen,
  onClose,
  onConfirm,
  title,
  description,
  confirmLabel = "Confirm",
  cancelLabel = "Cancel",
  variant = "default",
  requireConfirmationText,
}: ConfirmDialogProps) {
  const [loading, setLoading] = useState(false);
  const [confirmInput, setConfirmInput] = useState("");

  if (!isOpen) return null;

  const isConfirmDisabled =
    loading ||
    (requireConfirmationText !== undefined && confirmInput !== requireConfirmationText);

  const handleConfirm = async () => {
    try {
      setLoading(true);
      await onConfirm();
      onClose();
    } finally {
      setLoading(false);
      setConfirmInput("");
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

        <div className="flex items-start gap-4 mb-4">
          <div
            className={clsx(
              "flex h-10 w-10 shrink-0 items-center justify-center rounded-full",
              variant === "destructive"
                ? "bg-rose-500/15 text-rose-600 dark:text-rose-400"
                : "bg-primary/15 text-primary"
            )}
          >
            <AlertTriangle size={20} />
          </div>
          <div>
            <h3 className="text-lg font-semibold tracking-tight text-foreground">{title}</h3>
            <div className="text-sm text-muted-foreground mt-1">{description}</div>
          </div>
        </div>

        {requireConfirmationText && (
          <div className="my-4 space-y-1.5">
            <label className="text-xs font-medium text-muted-foreground">
              Type <span className="font-mono font-bold text-foreground">{requireConfirmationText}</span> to confirm:
            </label>
            <input
              type="text"
              value={confirmInput}
              onChange={(e) => setConfirmInput(e.target.value)}
              placeholder={requireConfirmationText}
              className="w-full h-9 px-3 rounded-md border bg-background text-sm font-mono focus:outline-none focus:ring-2 focus:ring-ring"
            />
          </div>
        )}

        <div className="flex items-center justify-end gap-3 mt-6">
          <Button
            variant="outline"
            onClick={onClose}
            disabled={loading}
          >
            {cancelLabel}
          </Button>
          <Button
            variant={variant === "destructive" ? "destructive" : "default"}
            onClick={handleConfirm}
            disabled={isConfirmDisabled}
          >
            {loading ? "Processing..." : confirmLabel}
          </Button>
        </div>
      </div>
    </div>
  );
}
