"use client";

import React from "react";
import { useParams } from "next/navigation";
import { PageHeader } from "@/components/ui/PageHeader";
import { ObjectExplorer } from "@/features/objects/ObjectExplorer";

export default function BucketDetailPage() {
  const params = useParams();
  const rawName = Array.isArray(params?.name) ? params.name[0] : params?.name || "";
  const bucketName = decodeURIComponent(rawName);

  return (
    <div className="space-y-6">
      <PageHeader
        title={bucketName}
        description={`Explore files and folders in "${bucketName}". Multi-cloud storage routing and policy evaluation are active.`}
      />

      <ObjectExplorer bucket={bucketName} />
    </div>
  );
}
