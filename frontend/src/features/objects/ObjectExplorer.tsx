"use client";

import React, { useState, useMemo } from "react";
import { useObjects } from "./hooks";
import { ObjectViewModel } from "./types";
import { ObjectTable } from "./ObjectTable";
import { ObjectBreadcrumbs } from "./ObjectBreadcrumbs";
import { DeleteObjectConfirm } from "./DeleteObjectConfirm";
import { UploadDropzone } from "@/features/uploads/UploadDropzone";
import { SearchBar } from "@/components/ui/SearchBar";
import { Button } from "@/components/ui/button";
import { RefreshCw, UploadCloud, Eye, EyeOff } from "lucide-react";

export interface ObjectExplorerProps {
  bucket: string;
}

export function ObjectExplorer({ bucket }: ObjectExplorerProps) {
  const [prefix, setPrefix] = useState("");
  const [searchQuery, setSearchQuery] = useState("");
  const [showDropzone, setShowDropzone] = useState(false);
  const [deletingObject, setDeletingObject] = useState<ObjectViewModel | null>(null);

  const { data: objects = [], isLoading, refetch } = useObjects(bucket, prefix);

  const filteredObjects = useMemo(() => {
    if (!searchQuery.trim()) return objects;
    const query = searchQuery.toLowerCase();
    return objects.filter((o) => o.name.toLowerCase().includes(query));
  }, [objects, searchQuery]);

  return (
    <div className="space-y-6">
      {/* Top Breadcrumb & Action bar */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-card p-4 rounded-xl border shadow-2xs">
        <ObjectBreadcrumbs
          bucket={bucket}
          prefix={prefix}
          onNavigate={(newPrefix) => {
            setPrefix(newPrefix);
            setSearchQuery("");
          }}
        />

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => refetch()}
            className="h-9 px-3"
            title="Refresh folder contents"
          >
            <RefreshCw size={14} className={isLoading ? "animate-spin" : ""} />
          </Button>

          <Button
            variant="outline"
            size="sm"
            onClick={() => setShowDropzone(!showDropzone)}
            className="h-9 gap-1.5 text-xs font-semibold"
          >
            {showDropzone ? <EyeOff size={15} /> : <Eye size={15} />}
            <span>{showDropzone ? "Hide Dropzone" : "Show Dropzone"}</span>
          </Button>

          <UploadDropzone
            bucket={bucket}
            prefix={prefix}
            compact
            onUploadSuccess={() => refetch()}
          />
        </div>
      </div>

      {/* Optional Large Dropzone */}
      {showDropzone && (
        <div className="animate-fade-in">
          <UploadDropzone
            bucket={bucket}
            prefix={prefix}
            onUploadSuccess={() => {
              refetch();
              setShowDropzone(false);
            }}
          />
        </div>
      )}

      {/* Filter & Search */}
      <div className="flex items-center justify-between gap-4">
        <SearchBar
          value={searchQuery}
          onChange={setSearchQuery}
          placeholder="Filter files in this folder..."
        />
      </div>

      {/* Objects Table */}
      <div className="animate-fade-in">
        <ObjectTable
          bucket={bucket}
          objects={filteredObjects}
          isLoading={isLoading}
          onSelectFolder={(folderKey) => {
            setPrefix(folderKey);
            setSearchQuery("");
          }}
          onDeleteObject={(item) => setDeletingObject(item)}
          emptyAction={
            <Button onClick={() => setShowDropzone(true)} className="gap-1.5 font-semibold">
              <UploadCloud size={16} />
              <span>Upload first file</span>
            </Button>
          }
        />
      </div>

      {/* Modals */}
      <DeleteObjectConfirm
        bucket={bucket}
        objectItem={deletingObject}
        onClose={() => setDeletingObject(null)}
      />
    </div>
  );
}
