import React from "react";
import { clsx } from "clsx";

export interface LoadingSkeletonProps {
  type?: "card" | "table" | "stats" | "line";
  count?: number;
  className?: string;
}

export function LoadingSkeleton({
  type = "line",
  count = 1,
  className,
}: LoadingSkeletonProps) {
  const renderSkeleton = (index: number) => {
    if (type === "card") {
      return (
        <div
          key={index}
          className={clsx(
            "rounded-xl border bg-card p-6 shadow-sm animate-pulse space-y-4",
            className
          )}
        >
          <div className="h-5 w-2/3 rounded bg-muted" />
          <div className="h-4 w-1/2 rounded bg-muted/70" />
          <div className="flex justify-between pt-4 border-t">
            <div className="h-4 w-16 rounded bg-muted" />
            <div className="h-4 w-16 rounded bg-muted" />
          </div>
        </div>
      );
    }

    if (type === "stats") {
      return (
        <div
          key={index}
          className={clsx(
            "rounded-xl border bg-card p-6 shadow-sm animate-pulse space-y-2",
            className
          )}
        >
          <div className="h-4 w-24 rounded bg-muted" />
          <div className="h-8 w-32 rounded bg-muted/80" />
          <div className="h-3 w-40 rounded bg-muted/60" />
        </div>
      );
    }

    if (type === "table") {
      return (
        <div key={index} className="w-full space-y-3 animate-pulse">
          <div className="h-10 w-full rounded bg-muted/60" />
          {Array.from({ length: 5 }).map((_, idx) => (
            <div
              key={idx}
              className="h-12 w-full rounded border bg-card/60 flex items-center px-4 gap-4"
            >
              <div className="h-4 w-1/4 rounded bg-muted" />
              <div className="h-4 w-1/6 rounded bg-muted" />
              <div className="h-4 w-1/6 rounded bg-muted" />
              <div className="h-4 w-1/6 rounded bg-muted" />
            </div>
          ))}
        </div>
      );
    }

    return (
      <div
        key={index}
        className={clsx("h-4 w-full rounded bg-muted animate-pulse", className)}
      />
    );
  };

  return (
    <>
      {Array.from({ length: count }).map((_, idx) => renderSkeleton(idx))}
    </>
  );
}
