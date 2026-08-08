"use client";

import React, { useState } from "react";
import { Laptop, Smartphone, Globe, Monitor, LogOut } from "lucide-react";
import { Button } from "@/components/ui/button";

interface Session {
  id: string;
  device: string;
  browser: string;
  os: string;
  ip: string;
  location: string;
  lastActive: string;
  isCurrent: boolean;
  type: "desktop" | "mobile" | "web";
}

const MOCK_SESSIONS: Session[] = [
  {
    id: "sess_1",
    device: "MacBook Pro M3",
    browser: "Chrome",
    os: "macOS",
    ip: "192.168.1.100",
    location: "San Francisco, CA, US",
    lastActive: "Active now",
    isCurrent: true,
    type: "desktop",
  },
  {
    id: "sess_2",
    device: "iPhone 15 Pro",
    browser: "Safari",
    os: "iOS",
    ip: "10.0.0.52",
    location: "San Francisco, CA, US",
    lastActive: "2 hours ago",
    isCurrent: false,
    type: "mobile",
  },
  {
    id: "sess_3",
    device: "Linux Workstation",
    browser: "Firefox",
    os: "Ubuntu",
    ip: "198.51.100.12",
    location: "New York, NY, US",
    lastActive: "3 days ago",
    isCurrent: false,
    type: "desktop",
  },
];

export default function SessionsPage() {
  const [sessions, setSessions] = useState(MOCK_SESSIONS);

  const handleRevoke = (id: string) => {
    setSessions(sessions.filter((s) => s.id !== id));
  };

  const handleRevokeAll = () => {
    setSessions(sessions.filter((s) => s.isCurrent));
  };

  const getIcon = (type: string) => {
    switch (type) {
      case "desktop": return <Laptop className="w-5 h-5 text-muted-foreground" />;
      case "mobile": return <Smartphone className="w-5 h-5 text-muted-foreground" />;
      default: return <Globe className="w-5 h-5 text-muted-foreground" />;
    }
  };

  return (
    <div className="space-y-8">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Active Sessions</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Manage your signed-in devices and sessions.
          </p>
        </div>
        <Button variant="destructive" onClick={handleRevokeAll}>
          <LogOut size={16} className="mr-2" />
          Revoke All Other Sessions
        </Button>
      </div>

      <div className="border rounded-xl bg-card overflow-hidden">
        <div className="grid gap-0 divide-y divide-border">
          {sessions.map((session) => (
            <div key={session.id} className="p-4 sm:p-6 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
              <div className="flex items-start gap-4">
                <div className="mt-1 p-2 bg-muted rounded-md shrink-0">
                  {getIcon(session.type)}
                </div>
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <p className="font-medium">{session.device}</p>
                    {session.isCurrent && (
                      <span className="px-2 py-0.5 text-[10px] uppercase tracking-wider font-semibold bg-emerald-500/15 text-emerald-500 rounded-full">
                        Current Device
                      </span>
                    )}
                  </div>
                  <p className="text-sm text-muted-foreground">
                    {session.browser} on {session.os} • {session.ip}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {session.location} • {session.lastActive}
                  </p>
                </div>
              </div>
              
              {!session.isCurrent && (
                <Button 
                  variant="outline" 
                  size="sm" 
                  className="w-full sm:w-auto"
                  onClick={() => handleRevoke(session.id)}
                >
                  Revoke
                </Button>
              )}
            </div>
          ))}

          {sessions.length === 0 && (
            <div className="p-8 text-center text-muted-foreground">
              No active sessions found.
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
