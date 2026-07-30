"use client";

import React from "react";
import { PageHeader } from "@/components/ui/PageHeader";
import { OverviewStats } from "@/features/dashboard/OverviewStats";
import { RecentBucketsCard } from "@/features/dashboard/RecentBucketsCard";

export default function DashboardOverviewPage() {
  return (
    <div className="space-y-6">
      <PageHeader
        title="Dashboard Overview"
        description="Monitor your multi-cloud storage infrastructure, buckets, and routing health."
      />

      <OverviewStats />

      <div className="pt-2">
        <RecentBucketsCard />
      </div>
    </div>
  );
}
