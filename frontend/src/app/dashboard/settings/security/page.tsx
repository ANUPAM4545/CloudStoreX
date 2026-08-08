"use client";

import React from "react";
import { Activity, Download, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export default function AuditLogPage() {
  const MOCK_LOGS = [
    { id: 1, action: "User Login", user: "anupam@example.com", resource: "Identity", ip: "192.168.1.100", date: "Just now" },
    { id: 2, action: "Revoked API Key", user: "anupam@example.com", resource: "API Key (key_3)", ip: "192.168.1.100", date: "2 mins ago" },
    { id: 3, action: "Updated Policy", user: "admin@cloudstorex.com", resource: "Policy (S3-Read)", ip: "10.0.0.52", date: "1 hour ago" },
    { id: 4, action: "Failed Login", user: "unknown", resource: "Identity", ip: "203.0.113.42", date: "3 hours ago" },
  ];

  return (
    <div className="space-y-8">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Audit Log</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Review security events and administrative actions across your organization.
          </p>
        </div>
        <Button variant="outline">
          <Download size={16} className="mr-2" />
          Export CSV
        </Button>
      </div>

      <div className="border rounded-xl bg-card overflow-hidden">
        <div className="p-4 border-b bg-muted/20">
          <div className="relative max-w-sm">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
            <Input placeholder="Search logs..." className="pl-9 h-10" />
          </div>
        </div>
        <table className="w-full text-sm text-left">
          <thead className="bg-muted/50 border-b text-muted-foreground">
            <tr>
              <th className="px-6 py-4 font-medium">Action</th>
              <th className="px-6 py-4 font-medium">User</th>
              <th className="px-6 py-4 font-medium">Resource</th>
              <th className="px-6 py-4 font-medium">IP Address</th>
              <th className="px-6 py-4 font-medium text-right">Timestamp</th>
            </tr>
          </thead>
          <tbody className="divide-y">
            {MOCK_LOGS.map((log) => (
              <tr key={log.id} className="hover:bg-muted/30 transition-colors">
                <td className="px-6 py-4 font-medium">{log.action}</td>
                <td className="px-6 py-4 text-muted-foreground">{log.user}</td>
                <td className="px-6 py-4 text-muted-foreground">{log.resource}</td>
                <td className="px-6 py-4 font-mono text-xs text-muted-foreground">{log.ip}</td>
                <td className="px-6 py-4 text-right text-muted-foreground">{log.date}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
