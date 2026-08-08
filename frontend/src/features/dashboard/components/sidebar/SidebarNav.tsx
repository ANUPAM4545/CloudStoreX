"use client";

import * as React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import { useSidebarStore } from "@/store/sidebar";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  LayoutDashboard,
  Database,
  FileText,
  Cloud,
  GitBranch,
  ChartColumn,
  Sparkles,
  Shield,
  FileSearch,
  Activity,
  Cpu,
  Copy,
  LifeBuoy,
  UserCircle,
  Building,
  Settings2,
} from "lucide-react";

interface NavGroup {
  label: string;
  items: NavItem[];
}

interface NavItem {
  title: string;
  href: string;
  icon: React.ElementType;
}

const NAV_GROUPS: NavGroup[] = [
  {
    label: "Platform",
    items: [
      { title: "Overview", href: "/dashboard/overview", icon: LayoutDashboard },
      { title: "Buckets", href: "/dashboard/buckets", icon: Database },
      { title: "Objects", href: "/dashboard/objects", icon: FileText },
      { title: "Providers", href: "/dashboard/providers", icon: Cloud },
      { title: "Policies", href: "/dashboard/policies", icon: GitBranch },
    ],
  },
  {
    label: "Intelligence",
    items: [
      { title: "AI Center", href: "/dashboard/ai", icon: Sparkles },
      { title: "Analytics", href: "/dashboard/analytics", icon: ChartColumn },
      { title: "Audit Logs", href: "/dashboard/audit", icon: FileSearch },
      { title: "Security", href: "/dashboard/security", icon: Shield },
    ],
  },
  {
    label: "Infrastructure",
    items: [
      { title: "Monitoring", href: "/dashboard/monitoring", icon: Activity },
      { title: "Jobs", href: "/dashboard/jobs", icon: Cpu },
      { title: "Replication", href: "/dashboard/replication", icon: Copy },
      { title: "Disaster Recovery", href: "/dashboard/dr", icon: LifeBuoy },
    ],
  },
  {
    label: "Administration",
    items: [
      { title: "Profile", href: "/dashboard/settings/profile", icon: UserCircle },
      { title: "Workspace", href: "/dashboard/settings/workspace", icon: Building },
      { title: "Configuration", href: "/dashboard/settings/configuration", icon: Settings2 },
    ],
  },
];

export function SidebarNav() {
  const pathname = usePathname();
  const { isCollapsed } = useSidebarStore();

  return (
    <div className="flex-1 overflow-y-auto overflow-x-hidden p-3 space-y-6 scrollbar-none">
      {NAV_GROUPS.map((group) => (
        <div key={group.label} className="space-y-1">
          <AnimatePresence initial={false}>
            {!isCollapsed && (
              <motion.div
                initial={{ opacity: 0, height: 0 }}
                animate={{ opacity: 1, height: "auto" }}
                exit={{ opacity: 0, height: 0 }}
                className="px-3 pb-1"
              >
                <h3 className="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground">
                  {group.label}
                </h3>
              </motion.div>
            )}
          </AnimatePresence>

          <nav className="space-y-1">
            {group.items.map((item) => {
              const isActive = pathname === item.href || pathname?.startsWith(`${item.href}/`);

              return (
                <Tooltip key={item.href}>
                  <TooltipTrigger
                    render={
                      <Link
                        href={item.href}
                        className={`
                          relative flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200 group
                          ${isCollapsed ? "justify-center" : "justify-start"}
                          ${isActive ? "text-foreground bg-primary/10" : "text-muted-foreground hover:bg-muted/50 hover:text-foreground"}
                        `}
                      />
                    }
                  >
                      {isActive && (
                        <motion.div
                          layoutId="activeNavIndicator"
                          className="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-5 bg-primary rounded-r-full"
                          initial={{ opacity: 0 }}
                          animate={{ opacity: 1 }}
                          transition={{ type: "spring", stiffness: 300, damping: 30 }}
                        />
                      )}
                      <item.icon size={18} className={`shrink-0 transition-colors ${isActive ? "text-primary" : "group-hover:text-foreground"}`} />
                      
                      <AnimatePresence initial={false}>
                        {!isCollapsed && (
                          <motion.span
                            initial={{ opacity: 0, width: 0 }}
                            animate={{ opacity: 1, width: "auto" }}
                            exit={{ opacity: 0, width: 0 }}
                            className="truncate"
                          >
                            {item.title}
                          </motion.span>
                        )}
                      </AnimatePresence>
                  </TooltipTrigger>
                  {isCollapsed && (
                    <TooltipContent side="right" className="flex items-center gap-2" sideOffset={12}>
                      <span className="font-semibold">{item.title}</span>
                    </TooltipContent>
                  )}
                </Tooltip>
              );
            })}
          </nav>
        </div>
      ))}
    </div>
  );
}
