"use client";

import { useEffect } from "react";
import { useRouter, usePathname } from "next/navigation";
import Link from "next/link";
import { useAuthStore } from "@/store/auth";
import { UploadProgress } from "@/features/uploads/UploadProgress";
import {
  Database,
  Folder,
  Settings,
  UserCircle,
  LogOut,
  LayoutDashboard,
  Cloud,
} from "lucide-react";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const { isAuthenticated, user, logout } = useAuthStore();

  useEffect(() => {
    if (!isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, router]);

  if (!isAuthenticated) {
    return null;
  }

  const isOverviewActive = pathname === "/dashboard/overview";
  const isBucketsActive =
    pathname === "/dashboard/buckets" || pathname?.startsWith("/dashboard/bucket/");

  return (
    <div className="flex h-screen overflow-hidden bg-background">
      {/* Sidebar */}
      <aside className="w-64 border-r bg-muted/30 flex-col hidden md:flex">
        <div className="p-6 flex items-center gap-2.5">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground shadow-xs">
            <Cloud size={18} />
          </div>
          <h2 className="text-xl font-bold tracking-tight text-foreground">
            CloudStoreX
          </h2>
        </div>

        <div className="px-4 py-2">
          <div className="w-full p-2.5 text-sm border rounded-lg bg-card shadow-2xs flex items-center justify-between">
            <div className="flex flex-col min-w-0">
              <span className="font-semibold text-xs text-foreground truncate">
                Default Workspace
              </span>
              <span className="text-[11px] text-muted-foreground">
                Policy: Standard
              </span>
            </div>
            <span className="text-xs text-primary font-medium bg-primary/10 px-2 py-0.5 rounded-full">
              MinIO
            </span>
          </div>
        </div>

        <nav className="flex-1 px-4 py-4 space-y-1.5">
          <Link
            href="/dashboard/overview"
            className={`flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
              isOverviewActive
                ? "bg-primary text-primary-foreground shadow-xs"
                : "text-muted-foreground hover:bg-muted/60 hover:text-foreground"
            }`}
          >
            <LayoutDashboard size={18} />
            <span>Overview</span>
          </Link>

          <Link
            href="/dashboard/buckets"
            className={`flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
              isBucketsActive
                ? "bg-primary text-primary-foreground shadow-xs"
                : "text-muted-foreground hover:bg-muted/60 hover:text-foreground"
            }`}
          >
            <Folder size={18} />
            <span>Buckets</span>
          </Link>

          <Link
            href="/dashboard/buckets"
            className="flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium text-muted-foreground hover:bg-muted/60 hover:text-foreground transition-colors"
          >
            <Database size={18} />
            <span>Providers</span>
          </Link>

          <Link
            href="/dashboard/overview"
            className="flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium text-muted-foreground hover:bg-muted/60 hover:text-foreground transition-colors"
          >
            <Settings size={18} />
            <span>Settings</span>
          </Link>
        </nav>

        <div className="p-4 border-t bg-muted/20">
          <div className="flex items-center gap-3 mb-3 px-2">
            <UserCircle size={32} className="text-muted-foreground shrink-0" />
            <div className="flex flex-col min-w-0">
              <span className="text-sm font-semibold truncate">
                {user?.fullName || "User"}
              </span>
              <span className="text-xs text-muted-foreground truncate" title={user?.email}>
                {user?.email || "user@example.com"}
              </span>
            </div>
          </div>
          <button
            onClick={() => {
              logout();
              router.push("/login");
            }}
            className="flex items-center gap-2 w-full px-3 py-2 text-sm font-medium text-destructive hover:bg-destructive/10 rounded-lg transition-colors"
          >
            <LogOut size={16} />
            <span>Sign Out</span>
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 flex flex-col min-w-0 overflow-hidden">
        {/* Mobile Header */}
        <header className="md:hidden h-14 border-b flex items-center justify-between px-4 bg-card">
          <div className="flex items-center gap-2">
            <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-primary text-primary-foreground">
              <Cloud size={16} />
            </div>
            <h2 className="text-base font-bold tracking-tight">CloudStoreX</h2>
          </div>
          <button
            onClick={() => {
              logout();
              router.push("/login");
            }}
            className="text-xs font-semibold text-destructive"
          >
            Sign Out
          </button>
        </header>

        <div className="flex-1 overflow-auto p-6">{children}</div>
      </main>

      {/* Global Upload Progress Manager */}
      <UploadProgress />
    </div>
  );
}
