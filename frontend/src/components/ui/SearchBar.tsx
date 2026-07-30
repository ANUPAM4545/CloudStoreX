import React from "react";
import { Search, X } from "lucide-react";
import { Input } from "@/components/ui/input";
import { clsx } from "clsx";

export interface SearchBarProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
}

export function SearchBar({
  value,
  onChange,
  placeholder = "Search...",
  className,
}: SearchBarProps) {
  return (
    <div className={clsx("relative flex items-center w-full max-w-sm", className)}>
      <Search
        size={16}
        className="absolute left-3 text-muted-foreground pointer-events-none"
      />
      <Input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        className="pl-9 pr-9 h-9 text-sm"
      />
      {value && (
        <button
          type="button"
          onClick={() => onChange("")}
          className="absolute right-2 p-1 text-muted-foreground hover:text-foreground rounded-full hover:bg-muted transition-colors"
          title="Clear search"
        >
          <X size={14} />
        </button>
      )}
    </div>
  );
}
