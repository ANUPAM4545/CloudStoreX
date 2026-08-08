"use client";

import React from "react";
import { Bot, Plus, Copy, Shield, MoreVertical } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function ServiceAccountsPage() {
  return (
    <div className="space-y-8">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Service Accounts</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Manage machine identities for automated workloads and CI/CD pipelines.
          </p>
        </div>
        <Button>
          <Plus size={16} className="mr-2" />
          Create Service Account
        </Button>
      </div>

      <div className="border rounded-xl bg-card overflow-hidden">
        <table className="w-full text-sm text-left">
          <thead className="bg-muted/50 border-b text-muted-foreground">
            <tr>
              <th className="px-6 py-4 font-medium">Service Account</th>
              <th className="px-6 py-4 font-medium">Role</th>
              <th className="px-6 py-4 font-medium">Status</th>
              <th className="px-6 py-4 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y">
            <tr className="hover:bg-muted/30 transition-colors group">
              <td className="px-6 py-4">
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-muted rounded-md text-foreground">
                    <Bot size={16} />
                  </div>
                  <div>
                    <p className="font-medium">Data Pipeline Worker</p>
                    <p className="text-xs text-muted-foreground font-mono mt-0.5">sa-pipeline@cloudstorex.local</p>
                  </div>
                </div>
              </td>
              <td className="px-6 py-4">
                <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-blue-500/10 text-blue-500 text-xs font-medium">
                  <Shield size={12} />
                  Storage Writer
                </span>
              </td>
              <td className="px-6 py-4">
                <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full border text-xs font-medium bg-emerald-500/10 text-emerald-500 border-emerald-500/20">
                  <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                  Active
                </span>
              </td>
              <td className="px-6 py-4 text-right">
                <Button variant="ghost" size="icon" className="text-muted-foreground hover:text-foreground">
                  <MoreVertical size={16} />
                </Button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  );
}
