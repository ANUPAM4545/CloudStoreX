"use client";

import React from "react";
import { ChevronRight, Folder, Home } from "lucide-react";

export interface ObjectBreadcrumbsProps {
  bucket: string;
  prefix: string;
  onNavigate: (prefix: string) => void;
}

export function ObjectBreadcrumbs({ bucket, prefix, onNavigate }: ObjectBreadcrumbsProps) {
  const parts = prefix.split("/").filter(Boolean);

  return (
    <nav aria-label="Breadcrumb" className="flex items-center gap-1 text-sm flex-wrap">
      <button
        type="button"
        onClick={() => onNavigate("")}
        className="flex items-center gap-1.5 font-semibold text-foreground hover:text-primary transition-colors px-1.5 py-1 rounded-md hover:bg-muted/50"
      >
        <Home size={15} />
        <span>{bucket}</span>
      </button>

      {parts.map((part, index) => {
        const path = parts.slice(0, index + 1).join("/") + "/";
        const isLast = index === parts.length - 1;

        return (
          <React.Fragment key={path}>
            <ChevronRight size={14} className="text-muted-foreground shrink-0" />
            <button
              type="button"
              onClick={() => !isLast && onNavigate(path)}
              disabled={isLast}
              className={`flex items-center gap-1.5 px-1.5 py-1 rounded-md transition-colors ${
                isLast
                  ? "font-bold text-primary cursor-default"
                  : "font-medium text-muted-foreground hover:text-foreground hover:bg-muted/50"
              }`}
            >
              <Folder size={14} className={isLast ? "text-primary" : "text-muted-foreground"} />
              <span>{part}</span>
            </button>
          </React.Fragment>
        );
      })}
    </nav>
  );
}
