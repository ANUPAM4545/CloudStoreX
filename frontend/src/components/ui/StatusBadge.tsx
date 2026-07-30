import React from "react";
import { clsx } from "clsx";

export type StatusType = "success" | "warning" | "error" | "info" | "neutral";

export interface StatusBadgeProps {
  status?: StatusType;
  label: string;
  dot?: boolean;
  className?: string;
}

export function StatusBadge({
  status = "success",
  label,
  dot = true,
  className,
}: StatusBadgeProps) {
  const statusStyles: Record<StatusType, string> = {
    success: "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-500/30",
    warning: "bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/30",
    error: "bg-rose-500/15 text-rose-700 dark:text-rose-400 border-rose-500/30",
    info: "bg-blue-500/15 text-blue-700 dark:text-blue-400 border-blue-500/30",
    neutral: "bg-muted text-muted-foreground border-border",
  };

  const dotStyles: Record<StatusType, string> = {
    success: "bg-emerald-500",
    warning: "bg-amber-500",
    error: "bg-rose-500",
    info: "bg-blue-500",
    neutral: "bg-muted-foreground",
  };

  return (
    <span
      className={clsx(
        "inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium border",
        statusStyles[status],
        className
      )}
    >
      {dot && (
        <span
          className={clsx("w-1.5 h-1.5 rounded-full animate-pulse", dotStyles[status])}
        />
      )}
      {label}
    </span>
  );
}
