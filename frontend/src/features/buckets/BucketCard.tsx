"use client";

import React from "react";
import Link from "next/link";
import { Folder, Trash2, ExternalLink, Calendar } from "lucide-react";
import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { BucketViewModel } from "./types";

export interface BucketCardProps {
  bucket: BucketViewModel;
  onDelete: (bucket: BucketViewModel) => void;
}

export function BucketCard({ bucket, onDelete }: BucketCardProps) {
  return (
    <Card className="group relative overflow-hidden transition-all hover:shadow-md hover:border-primary/50 bg-card flex flex-col justify-between">
      <div>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-colors">
              <Folder size={20} />
            </div>
            <div className="flex flex-col min-w-0">
              <Link
                href={`/dashboard/bucket/${encodeURIComponent(bucket.name)}`}
                className="font-semibold text-foreground hover:text-primary truncate transition-colors text-base"
              >
                {bucket.name}
              </Link>
              <span className="text-xs text-muted-foreground">Standard Storage</span>
            </div>
          </div>
        </CardHeader>

        <CardContent className="pt-3 pb-4">
          <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
            <Calendar size={13} />
            <span>Created {bucket.createdAtFormatted}</span>
          </div>
        </CardContent>
      </div>

      <CardFooter className="flex items-center justify-between border-t pt-3 pb-3 bg-muted/20">
        <Link href={`/dashboard/bucket/${encodeURIComponent(bucket.name)}`}>
          <Button variant="ghost" size="sm" className="h-8 gap-1.5 text-xs font-medium">
            <span>Explore</span>
            <ExternalLink size={13} />
          </Button>
        </Link>
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onDelete(bucket)}
          className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive hover:bg-destructive/10"
          title="Delete bucket"
        >
          <Trash2 size={15} />
        </Button>
      </CardFooter>
    </Card>
  );
}
