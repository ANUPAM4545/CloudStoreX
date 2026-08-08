"use client";

import * as React from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Search } from "lucide-react";
import { useSidebarStore } from "@/store/sidebar";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

export function SidebarSearch() {
  const { isCollapsed } = useSidebarStore();

  return (
    <div className="px-3 pt-3 pb-2">
      <Tooltip>
        <TooltipTrigger
          render={
            <button
              onClick={() => {
                // Trigger command palette later
                console.log("Open command palette");
              }}
              className={`
                w-full flex items-center justify-between rounded-lg border border-input bg-background/50 text-muted-foreground
                hover:bg-muted/80 hover:text-foreground transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/50
                ${isCollapsed ? "h-10 px-0 justify-center" : "h-9 px-3"}
              `}
            />
          }
        >
            <div className="flex items-center gap-2">
              <Search size={16} className="shrink-0" />
              <AnimatePresence initial={false}>
                {!isCollapsed && (
                  <motion.span
                    initial={{ opacity: 0, width: 0 }}
                    animate={{ opacity: 1, width: "auto" }}
                    exit={{ opacity: 0, width: 0 }}
                    className="text-sm font-medium whitespace-nowrap overflow-hidden"
                  >
                    Search...
                  </motion.span>
                )}
              </AnimatePresence>
            </div>
            
            <AnimatePresence initial={false}>
              {!isCollapsed && (
                <motion.kbd
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  exit={{ opacity: 0 }}
                  className="pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100"
                >
                  <span className="text-xs">⌘</span>K
                </motion.kbd>
              )}
            </AnimatePresence>
        </TooltipTrigger>
        {isCollapsed && (
          <TooltipContent side="right" sideOffset={12}>
            <div className="flex items-center gap-2">
              <span className="font-semibold">Search</span>
              <kbd className="inline-flex h-4 items-center rounded border bg-muted px-1 text-[10px] font-medium">⌘K</kbd>
            </div>
          </TooltipContent>
        )}
      </Tooltip>
    </div>
  );
}
