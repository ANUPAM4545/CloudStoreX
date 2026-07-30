"use client";

import React, { useState, useMemo } from "react";
import { useBuckets } from "./hooks";
import { BucketViewModel } from "./types";
import { BucketCard } from "./BucketCard";
import { BucketTable } from "./BucketTable";
import { CreateBucketModal } from "./CreateBucketModal";
import { DeleteBucketConfirm } from "./DeleteBucketConfirm";
import { SearchBar } from "@/components/ui/SearchBar";
import { Button } from "@/components/ui/button";
import { EmptyState } from "@/components/ui/EmptyState";
import { LoadingSkeleton } from "@/components/ui/LoadingSkeleton";
import { FolderPlus, LayoutGrid, List, RefreshCw } from "lucide-react";

export function BucketList() {
  const { data: buckets = [], isLoading, refetch } = useBuckets();
  const [searchQuery, setSearchQuery] = useState("");
  const [viewMode, setViewMode] = useState<"grid" | "table">("grid");
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [deletingBucket, setDeletingBucket] = useState<BucketViewModel | null>(null);

  const filteredBuckets = useMemo(() => {
    if (!searchQuery.trim()) return buckets;
    const query = searchQuery.toLowerCase();
    return buckets.filter((b) => b.name.toLowerCase().includes(query));
  }, [buckets, searchQuery]);

  return (
    <div className="space-y-6">
      {/* Action Bar */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <SearchBar
          value={searchQuery}
          onChange={setSearchQuery}
          placeholder="Search storage buckets..."
        />

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => refetch()}
            className="h-9 px-3"
            title="Refresh buckets list"
          >
            <RefreshCw size={14} className={isLoading ? "animate-spin" : ""} />
          </Button>

          <div className="flex items-center rounded-lg border bg-muted/40 p-0.5">
            <button
              type="button"
              onClick={() => setViewMode("grid")}
              className={`p-1.5 rounded-md text-xs transition-colors ${
                viewMode === "grid"
                  ? "bg-background text-foreground shadow-xs font-semibold"
                  : "text-muted-foreground hover:text-foreground"
              }`}
              title="Grid view"
            >
              <LayoutGrid size={15} />
            </button>
            <button
              type="button"
              onClick={() => setViewMode("table")}
              className={`p-1.5 rounded-md text-xs transition-colors ${
                viewMode === "table"
                  ? "bg-background text-foreground shadow-xs font-semibold"
                  : "text-muted-foreground hover:text-foreground"
              }`}
              title="Table view"
            >
              <List size={15} />
            </button>
          </div>

          <Button
            onClick={() => setIsCreateModalOpen(true)}
            className="h-9 gap-1.5 font-semibold"
          >
            <FolderPlus size={16} />
            <span>New Bucket</span>
          </Button>
        </div>
      </div>

      {/* Content Area */}
      {isLoading ? (
        viewMode === "grid" ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <LoadingSkeleton type="card" count={6} />
          </div>
        ) : (
          <LoadingSkeleton type="table" />
        )
      ) : filteredBuckets.length === 0 ? (
        <EmptyState
          icon={<FolderPlus size={32} />}
          title={searchQuery ? "No matching buckets" : "No storage buckets yet"}
          description={
            searchQuery
              ? `No buckets match "${searchQuery}". Try a different filter.`
              : "Get started by creating your first cloud storage bucket. CloudStoreX will manage policy routing automatically."
          }
          action={
            !searchQuery ? (
              <Button onClick={() => setIsCreateModalOpen(true)} className="gap-1.5 font-semibold">
                <FolderPlus size={16} />
                <span>Create your first bucket</span>
              </Button>
            ) : undefined
          }
        />
      ) : viewMode === "grid" ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 animate-fade-in">
          {filteredBuckets.map((bucket) => (
            <BucketCard
              key={bucket.name}
              bucket={bucket}
              onDelete={(b) => setDeletingBucket(b)}
            />
          ))}
        </div>
      ) : (
        <div className="animate-fade-in">
          <BucketTable
            buckets={filteredBuckets}
            onDelete={(b) => setDeletingBucket(b)}
            emptyAction={
              <Button onClick={() => setIsCreateModalOpen(true)} className="gap-1.5">
                <FolderPlus size={16} />
                <span>Create Bucket</span>
              </Button>
            }
          />
        </div>
      )}

      {/* Modals */}
      <CreateBucketModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
      />

      <DeleteBucketConfirm
        bucket={deletingBucket}
        onClose={() => setDeletingBucket(null)}
      />
    </div>
  );
}
