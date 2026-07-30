"use client";

import React from "react";
import { PageHeader } from "@/components/ui/PageHeader";
import { BucketList } from "@/features/buckets/BucketList";

export default function BucketsPage() {
  return (
    <div className="space-y-6">
      <PageHeader
        title="Storage Buckets"
        description="Create, filter, and inspect your cloud storage containers across providers."
      />

      <BucketList />
    </div>
  );
}
