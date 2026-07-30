"use client";

import React from "react";
import { ColumnDef } from "@tanstack/react-table";
import {
  Folder,
  FileText,
  FileCode,
  Image as ImageIcon,
  FileArchive,
  File,
  Download,
  Trash2,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { DataTable, SortableHeader } from "@/components/ui/DataTable";
import { ObjectViewModel } from "./types";
import { objectsApi } from "./api";

export interface ObjectTableProps {
  bucket: string;
  objects: ObjectViewModel[];
  isLoading?: boolean;
  onSelectFolder: (folderKey: string) => void;
  onDeleteObject: (objectItem: ObjectViewModel) => void;
  emptyAction?: React.ReactNode;
}

function getFileIcon(item: ObjectViewModel) {
  if (item.isFolder) return <Folder size={18} className="text-primary fill-primary/20" />;
  const name = item.name.toLowerCase();
  if (name.endsWith(".jpg") || name.endsWith(".png") || name.endsWith(".svg") || name.endsWith(".webp") || name.endsWith(".gif")) {
    return <ImageIcon size={18} className="text-blue-500" />;
  }
  if (name.endsWith(".zip") || name.endsWith(".tar") || name.endsWith(".gz") || name.endsWith(".7z")) {
    return <FileArchive size={18} className="text-amber-500" />;
  }
  if (name.endsWith(".json") || name.endsWith(".ts") || name.endsWith(".js") || name.endsWith(".go") || name.endsWith(".py") || name.endsWith(".html") || name.endsWith(".css")) {
    return <FileCode size={18} className="text-purple-500" />;
  }
  if (name.endsWith(".txt") || name.endsWith(".md") || name.endsWith(".pdf") || name.endsWith(".doc")) {
    return <FileText size={18} className="text-emerald-500" />;
  }
  return <File size={18} className="text-muted-foreground" />;
}

export function ObjectTable({
  bucket,
  objects,
  isLoading,
  onSelectFolder,
  onDeleteObject,
  emptyAction,
}: ObjectTableProps) {
  const columns: ColumnDef<ObjectViewModel>[] = [
    {
      accessorKey: "name",
      header: ({ column }) => <SortableHeader column={column} label="Name" />,
      cell: ({ row }) => {
        const item = row.original;
        return (
          <div className="flex items-center gap-3">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-muted/40 shrink-0">
              {getFileIcon(item)}
            </div>
            {item.isFolder ? (
              <button
                type="button"
                onClick={() => onSelectFolder(item.key)}
                className="font-semibold text-foreground hover:text-primary transition-colors text-left truncate max-w-sm"
              >
                {item.name}/
              </button>
            ) : (
              <span className="font-medium text-foreground truncate max-w-sm" title={item.key}>
                {item.name}
              </span>
            )}
          </div>
        );
      },
    },
    {
      accessorKey: "sizeRaw",
      header: ({ column }) => <SortableHeader column={column} label="Size" />,
      cell: ({ row }) => (
        <span className="text-muted-foreground font-mono text-xs">
          {row.original.sizeFormatted}
        </span>
      ),
    },
    {
      accessorKey: "lastModifiedFormatted",
      header: ({ column }) => <SortableHeader column={column} label="Last Modified" />,
      cell: ({ row }) => (
        <span className="text-muted-foreground text-xs">
          {row.original.lastModifiedFormatted}
        </span>
      ),
    },
    {
      id: "actions",
      header: () => <span className="text-right block">Actions</span>,
      cell: ({ row }) => {
        const item = row.original;
        if (item.isFolder) {
          return (
            <div className="flex justify-end">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => onSelectFolder(item.key)}
                className="h-8 text-xs font-semibold text-primary"
              >
                Open
              </Button>
            </div>
          );
        }

        return (
          <div className="flex items-center justify-end gap-1.5">
            <Button
              variant="outline"
              size="sm"
              onClick={() => objectsApi.downloadObject(bucket, item.key)}
              className="h-8 gap-1 text-xs"
              title="Download file"
            >
              <Download size={13} />
              <span>Download</span>
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onDeleteObject(item)}
              className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive hover:bg-destructive/10"
              title="Delete object"
            >
              <Trash2 size={14} />
            </Button>
          </div>
        );
      },
    },
  ];

  return (
    <DataTable
      columns={columns}
      data={objects}
      isLoading={isLoading}
      emptyTitle="No objects in this bucket folder"
      emptyDescription="Upload files to begin storing data in this location."
      emptyAction={emptyAction}
    />
  );
}
