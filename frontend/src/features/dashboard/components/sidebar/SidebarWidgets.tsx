"use client";

import * as React from "react";
import { motion, AnimatePresence } from "framer-motion";
import { useSidebarStore } from "@/store/sidebar";
import { Progress } from "@/components/ui/progress";

export function SidebarWidgets() {
  const { isCollapsed } = useSidebarStore();
  const [progress, setProgress] = React.useState(0);

  // Animate progress on mount
  React.useEffect(() => {
    const timer = setTimeout(() => setProgress(35), 500);
    return () => clearTimeout(timer);
  }, []);

  if (isCollapsed) return null;

  return (
    <motion.div 
      initial={{ opacity: 0, height: 0 }}
      animate={{ opacity: 1, height: "auto" }}
      exit={{ opacity: 0, height: 0 }}
      className="px-4 py-4 space-y-5 overflow-hidden"
    >
      {/* Provider Status */}
      <div className="space-y-2">
        <h4 className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">
          Providers
        </h4>
        <div className="space-y-1.5">
          <div className="flex items-center justify-between text-xs">
            <span className="font-medium">AWS</span>
            <div className="flex items-center gap-1.5">
              <span className="text-[10px] text-emerald-500 font-medium">Healthy</span>
              <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
            </div>
          </div>
          <div className="flex items-center justify-between text-xs">
            <span className="font-medium">MinIO</span>
            <div className="flex items-center gap-1.5">
              <span className="text-[10px] text-emerald-500 font-medium">Healthy</span>
              <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
            </div>
          </div>
          <div className="flex items-center justify-between text-xs">
            <span className="font-medium text-muted-foreground">Azure</span>
            <div className="flex items-center gap-1.5">
              <span className="text-[10px] text-muted-foreground">Coming Soon</span>
              <div className="w-1 h-1 rounded-full bg-muted-foreground/50" />
            </div>
          </div>
        </div>
      </div>

      {/* Storage Usage */}
      <div className="space-y-2.5 p-3 bg-muted/30 rounded-xl border border-border/50">
        <div className="flex justify-between items-end">
          <span className="text-xs font-semibold">Storage</span>
          <span className="text-[10px] font-medium text-muted-foreground">35 GB / 100 GB</span>
        </div>
        <Progress value={progress} className="h-1.5 bg-background" />
        <div className="text-[10px] text-muted-foreground flex items-center justify-between">
          <span>Enterprise Plan</span>
          <span className="text-primary hover:underline cursor-pointer">Upgrade</span>
        </div>
      </div>
    </motion.div>
  );
}
