"use client";

import * as React from "react";
import { motion } from "framer-motion";
import { PanelLeftClose, PanelLeftOpen } from "lucide-react";
import { useSidebarStore } from "@/store/sidebar";
import { WorkspaceSwitcher } from "./WorkspaceSwitcher";
import { SidebarSearch } from "./SidebarSearch";
import { SidebarNav } from "./SidebarNav";
import { SidebarWidgets } from "./SidebarWidgets";
import { SidebarProfile } from "./SidebarProfile";
import { CommandPalette } from "./CommandPalette";

export function Sidebar() {
  const { isCollapsed, toggleCollapse } = useSidebarStore();

  // Handle hydration mismatch safely
  const [mounted, setMounted] = React.useState(false);
  React.useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    // Render a skeleton of the expanded sidebar on server to prevent flash
    return (
      <aside className="w-[260px] h-screen border-r bg-muted/20 flex flex-col hidden md:flex shrink-0">
        <div className="h-16" />
      </aside>
    );
  }

  return (
    <motion.aside
      initial={false}
      animate={{ 
        width: isCollapsed ? 70 : 260 
      }}
      transition={{ 
        type: "spring", 
        stiffness: 300, 
        damping: 30 
      }}
      className="h-screen border-r bg-muted/20 flex flex-col hidden md:flex relative shrink-0 z-20 group/sidebar overflow-visible"
    >
      {/* Collapse Toggle Button (Floating) */}
      <button
        onClick={toggleCollapse}
        className="absolute -right-3 top-6 flex h-6 w-6 items-center justify-center rounded-full border bg-background text-muted-foreground shadow-sm hover:text-foreground hover:bg-muted focus-visible:outline-none focus-visible:ring-2 z-50 opacity-0 group-hover/sidebar:opacity-100 transition-opacity"
      >
        {isCollapsed ? <PanelLeftOpen size={12} /> : <PanelLeftClose size={12} />}
      </button>

      {/* Top Section */}
      <div className="p-3 border-b border-border/40">
        <WorkspaceSwitcher />
      </div>

      <SidebarSearch />

      {/* Main Navigation */}
      <SidebarNav />

      {/* Bottom Section */}
      <div className="mt-auto flex flex-col">
        <SidebarWidgets />
        <SidebarProfile />
      </div>

      {/* Global Command Palette */}
      <CommandPalette />
    </motion.aside>
  );
}
