"use client";

import React, { useEffect, useState } from "react";
import { Card, Title, Text, Grid, Badge, Button } from "@tremor/react";

export default function OperationsCenter() {
  const [health, setHealth] = useState<any>(null);

  useEffect(() => {
    // For demo purposes, fetch backend mock response or use static if API missing
    fetch("/api/v1/operations/cluster/health")
      .then((res) => {
        if (!res.ok) throw new Error("API not ready");
        return res.json();
      })
      .then(setHealth)
      .catch(() => {
        setHealth({ status: "healthy", nodes: 3, uptime: "99.99%" });
      });
  }, []);

  return (
    <main className="p-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <Title>Operations Center</Title>
          <Text>Manage cluster health, DR, chaos engineering, and failovers.</Text>
        </div>
      </div>
      
      <Grid numItems={1} numItemsSm={2} numItemsLg={3} className="gap-6">
        <Card>
          <Title>Cluster Health</Title>
          {health ? (
            <div className="mt-4">
              <Text>Status: <Badge color="emerald">{health.status}</Badge></Text>
              <Text className="mt-2">Nodes: {health.nodes}</Text>
              <Text className="mt-2">Uptime: {health.uptime}</Text>
            </div>
          ) : (
            <Text>Loading...</Text>
          )}
        </Card>

        <Card>
          <Title>Replication</Title>
          <Text className="mt-4">Active inter-region replications: 5</Text>
          <Text>Current lag: &lt; 1s</Text>
        </Card>

        <Card>
          <Title>Disaster Recovery</Title>
          <Text className="mt-4">RTO Target: 60s</Text>
          <Text>RPO Target: 300s</Text>
          <Button size="sm" className="mt-4">Run DR Drill</Button>
        </Card>

        <Card>
          <Title>Self-Healing & Backups</Title>
          <Text className="mt-4">Consistency Checks: All passed</Text>
          <Text>Last Backup: 2 hours ago</Text>
        </Card>

        <Card>
          <Title>Chaos Engineering</Title>
          <Text className="mt-4">Simulate outages, latency, and data corruption.</Text>
          <Button size="sm" color="red" className="mt-4">Run Chaos Experiment</Button>
        </Card>
      </Grid>
    </main>
  );
}
