"use client";

import React from "react";
import Link from "next/link";
import { ColumnDef } from "@tanstack/react-table";
import { Folder, Trash2, ExternalLink } from "lucide-react";
import { Button } from "@/components/ui/button";
import { DataTable, SortableHeader } from "@/components/ui/DataTable";
import { BucketViewModel } from "./types";

export interface BucketTableProps {
  buckets: BucketViewModel[];
  isLoading?: boolean;
  onDelete: (bucket: BucketViewModel) => void;
  emptyAction?: React.ReactNode;
}

export function BucketTable({ buckets, isLoading, onDelete, emptyAction }: BucketTableProps) {
  const columns: ColumnDef<BucketViewModel>[] = [
    {
      accessorKey: "name",
      header: ({ column }) => <SortableHeader column={column} label="Bucket Name" />,
      cell: ({ row }) => {
        const bucket = row.original;
        return (
          <div className="flex items-center gap-3">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10 text-primary">
              <Folder size={16} />
            </div>
            <Link
              href={`/dashboard/bucket/${encodeURIComponent(bucket.name)}`}
              className="font-semibold text-foreground hover:text-primary transition-colors"
            >
              {bucket.name}
            </Link>
          </div>
        );
      },
    },
    {
      accessorKey: "createdAtRaw",
      header: ({ column }) => <SortableHeader column={column} label="Created At" />,
      cell: ({ row }) => (
        <span className="text-muted-foreground">{row.original.createdAtFormatted}</span>
      ),
    },
    {
      id: "actions",
      header: () => <span className="text-right block">Actions</span>,
      cell: ({ row }) => {
        const bucket = row.original;
        return (
          <div className="flex items-center justify-end gap-2">
            <Link href={`/dashboard/bucket/${encodeURIComponent(bucket.name)}`}>
              <Button variant="outline" size="sm" className="h-8 gap-1.5 text-xs">
                <span>Explore</span>
                <ExternalLink size={12} />
              </Button>
            </Link>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onDelete(bucket)}
              className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive hover:bg-destructive/10"
              title="Delete bucket"
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
      data={buckets}
      isLoading={isLoading}
      emptyTitle="No buckets found"
      emptyDescription="No storage buckets match your criteria. Create your first bucket to begin."
      emptyAction={emptyAction}
    />
  );
}
