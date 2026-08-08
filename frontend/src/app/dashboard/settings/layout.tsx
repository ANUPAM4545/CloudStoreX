"use client";

import React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { Shield, Key, Fingerprint, Users, Laptop, Activity, Lock } from "lucide-react";

const SETTINGS_NAV = [
  { name: "Sessions", href: "/dashboard/settings/sessions", icon: <Laptop size={18} /> },
  { name: "API Keys", href: "/dashboard/settings/api-keys", icon: <Key size={18} /> },
  { name: "Service Accounts", href: "/dashboard/settings/service-accounts", icon: <Lock size={18} /> },
  { name: "MFA & Passkeys", href: "/dashboard/settings/mfa", icon: <Fingerprint size={18} /> },
  { name: "Enterprise SSO", href: "/dashboard/settings/sso", icon: <Shield size={18} /> },
  { name: "Organization", href: "/dashboard/settings/team", icon: <Users size={18} /> },
  { name: "Audit Log", href: "/dashboard/settings/security", icon: <Activity size={18} /> },
];

export default function SettingsLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();

  return (
    <div className="flex-1 flex flex-col sm:flex-row h-full overflow-hidden">
      {/* Settings Sidebar */}
      <aside className="w-full sm:w-64 border-r bg-muted/20 shrink-0 overflow-y-auto">
        <div className="p-6">
          <h2 className="text-lg font-semibold tracking-tight mb-4">Security Settings</h2>
          <nav className="space-y-1">
            {SETTINGS_NAV.map((item) => {
              const isActive = pathname === item.href;
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors ${
                    isActive 
                      ? "bg-primary/10 text-primary" 
                      : "text-muted-foreground hover:bg-muted hover:text-foreground"
                  }`}
                >
                  <span className={isActive ? "text-primary" : "text-muted-foreground"}>
                    {item.icon}
                  </span>
                  {item.name}
                </Link>
              );
            })}
          </nav>
        </div>
      </aside>

      {/* Main Content Area */}
      <main className="flex-1 overflow-y-auto p-6 md:p-10">
        <div className="max-w-4xl mx-auto">
          {children}
        </div>
      </main>
    </div>
  );
}
