"use client";

import { useEffect } from "react";
import { useRouter, usePathname } from "next/navigation";
import Link from "next/link";
import { useAuthStore } from "@/store/auth";
import { UploadProgress } from "@/features/uploads/UploadProgress";
import { Cloud } from "lucide-react";
import { Sidebar } from "@/features/dashboard/components/sidebar/Sidebar";

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

  // Layout active state is handled purely inside SidebarNav.tsx now

  return (
    <div className="flex h-screen overflow-hidden bg-background">
      {/* Enhanced Enterprise Sidebar */}
      <Sidebar />

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
