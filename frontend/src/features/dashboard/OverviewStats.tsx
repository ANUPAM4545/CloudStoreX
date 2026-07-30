"use client";

import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { Folder, HardDrive, Shield, Zap } from "lucide-react";
import { useBuckets } from "@/features/buckets/hooks";

export function OverviewStats() {
  const { data: buckets = [], isLoading } = useBuckets();

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {/* Total Buckets Card */}
      <Card className="bg-card hover:shadow-sm transition-shadow">
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
            Total Buckets
          </CardTitle>
          <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center text-primary">
            <Folder size={16} />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold tracking-tight text-foreground">
            {isLoading ? "…" : buckets.length}
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            Active storage containers
          </p>
        </CardContent>
      </Card>

      {/* Provider Status Card */}
      <Card className="bg-card hover:shadow-sm transition-shadow">
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
            Storage Provider
          </CardTitle>
          <div className="h-8 w-8 rounded-lg bg-emerald-500/10 flex items-center justify-center text-emerald-600 dark:text-emerald-400">
            <HardDrive size={16} />
          </div>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-2">
            <StatusBadge status="success" label="MinIO • Connected" />
          </div>
          <p className="text-xs text-muted-foreground mt-1.5">
            S3-compatible local cluster
          </p>
        </CardContent>
      </Card>

      {/* Policy Card */}
      <Card className="bg-card hover:shadow-sm transition-shadow">
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
            Routing Policy
          </CardTitle>
          <div className="h-8 w-8 rounded-lg bg-blue-500/10 flex items-center justify-center text-blue-600 dark:text-blue-400">
            <Shield size={16} />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-lg font-semibold tracking-tight text-foreground">
            Standard Default
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            Zero-trust provider isolation
          </p>
        </CardContent>
      </Card>

      {/* Max Upload Size Card */}
      <Card className="bg-card hover:shadow-sm transition-shadow">
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
            Max Object Size
          </CardTitle>
          <div className="h-8 w-8 rounded-lg bg-amber-500/10 flex items-center justify-center text-amber-600 dark:text-amber-400">
            <Zap size={16} />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold tracking-tight text-foreground">
            100 MB
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            Streaming multipart pipeline
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
