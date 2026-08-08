"use client";

import React from "react";
import { Users, Mail, MoreVertical } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function TeamPage() {
  return (
    <div className="space-y-8">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Organization Members</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Manage who has access to this CloudStoreX organization.
          </p>
        </div>
        <Button>
          <Mail size={16} className="mr-2" />
          Invite Members
        </Button>
      </div>

      <div className="border rounded-xl bg-card overflow-hidden">
        <table className="w-full text-sm text-left">
          <thead className="bg-muted/50 border-b text-muted-foreground">
            <tr>
              <th className="px-6 py-4 font-medium">Member</th>
              <th className="px-6 py-4 font-medium">Role</th>
              <th className="px-6 py-4 font-medium">MFA Status</th>
              <th className="px-6 py-4 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y">
            <tr className="hover:bg-muted/30 transition-colors group">
              <td className="px-6 py-4">
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-full bg-primary/20 text-primary flex items-center justify-center font-bold">
                    A
                  </div>
                  <div>
                    <p className="font-medium">Anupam (You)</p>
                    <p className="text-xs text-muted-foreground">anupam@example.com</p>
                  </div>
                </div>
              </td>
              <td className="px-6 py-4">
                <span className="font-medium">Owner</span>
              </td>
              <td className="px-6 py-4">
                <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full border text-xs font-medium bg-emerald-500/10 text-emerald-500 border-emerald-500/20">
                  <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                  Enabled
                </span>
              </td>
              <td className="px-6 py-4 text-right">
                <Button variant="ghost" size="icon" className="text-muted-foreground hover:text-foreground" disabled>
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
