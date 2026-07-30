"use client";

import React, { useState } from "react";
import Link from "next/link";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Folder, FolderPlus, ArrowRight, Calendar } from "lucide-react";
import { useBuckets } from "@/features/buckets/hooks";
import { CreateBucketModal } from "@/features/buckets/CreateBucketModal";
import { LoadingSkeleton } from "@/components/ui/LoadingSkeleton";

export function RecentBucketsCard() {
  const { data: buckets = [], isLoading } = useBuckets();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  const recentBuckets = buckets.slice(0, 5);

  return (
    <>
      <Card className="bg-card shadow-xs">
        <CardHeader className="flex flex-row items-center justify-between pb-3">
          <div>
            <CardTitle className="text-lg font-bold tracking-tight">Recent Storage Buckets</CardTitle>
            <CardDescription className="text-xs">
              Quickly jump into your recently created or modified buckets.
            </CardDescription>
          </div>
          <div className="flex items-center gap-2">
            <Button
              size="sm"
              variant="outline"
              onClick={() => setIsCreateModalOpen(true)}
              className="h-8 gap-1.5 text-xs font-semibold"
            >
              <FolderPlus size={14} />
              <span>New Bucket</span>
            </Button>
            {buckets.length > 0 && (
              <Link href="/dashboard/buckets">
                <Button size="sm" variant="ghost" className="h-8 gap-1 text-xs">
                  <span>View All</span>
                  <ArrowRight size={14} />
                </Button>
              </Link>
            )}
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <LoadingSkeleton type="table" />
          ) : recentBuckets.length === 0 ? (
            <div className="flex flex-col items-center justify-center p-8 text-center border rounded-xl border-dashed bg-muted/20">
              <FolderPlus size={28} className="text-muted-foreground mb-2" />
              <p className="text-sm font-semibold text-foreground">No buckets created yet</p>
              <p className="text-xs text-muted-foreground max-w-xs mt-1 mb-4">
                Create a bucket to store and organize your files across cloud providers.
              </p>
              <Button size="sm" onClick={() => setIsCreateModalOpen(true)} className="gap-1.5">
                <FolderPlus size={14} />
                <span>Create First Bucket</span>
              </Button>
            </div>
          ) : (
            <div className="divide-y border rounded-xl overflow-hidden">
              {recentBuckets.map((bucket) => (
                <div
                  key={bucket.name}
                  className="flex items-center justify-between p-3.5 hover:bg-muted/40 transition-colors"
                >
                  <div className="flex items-center gap-3 min-w-0">
                    <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10 text-primary">
                      <Folder size={18} />
                    </div>
                    <div className="min-w-0">
                      <Link
                        href={`/dashboard/bucket/${encodeURIComponent(bucket.name)}`}
                        className="font-semibold text-sm text-foreground hover:text-primary transition-colors block truncate"
                      >
                        {bucket.name}
                      </Link>
                      <div className="flex items-center gap-1.5 text-xs text-muted-foreground mt-0.5">
                        <Calendar size={12} />
                        <span>{bucket.createdAtFormatted}</span>
                      </div>
                    </div>
                  </div>
                  <Link href={`/dashboard/bucket/${encodeURIComponent(bucket.name)}`}>
                    <Button variant="outline" size="sm" className="h-8 text-xs gap-1">
                      <span>Explore</span>
                      <ArrowRight size={13} />
                    </Button>
                  </Link>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      <CreateBucketModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
      />
    </>
  );
}
