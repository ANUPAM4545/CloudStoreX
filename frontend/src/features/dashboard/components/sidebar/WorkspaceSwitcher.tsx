"use client";

import * as React from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Check, ChevronsUpDown, Cloud, Plus, Settings } from "lucide-react";
import { useSidebarStore } from "@/store/sidebar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export function WorkspaceSwitcher() {
  const { isCollapsed } = useSidebarStore();
  const [open, setOpen] = React.useState(false);

  return (
    <DropdownMenu open={open} onOpenChange={setOpen}>
      <DropdownMenuTrigger
        render={
          <button 
            className={`
              w-full flex items-center justify-between rounded-xl px-2 py-2 
              hover:bg-muted/50 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/50
              ${isCollapsed ? 'justify-center' : ''}
            `}
          />
        }
      >
          <div className="flex items-center gap-3 min-w-0">
            <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-blue-600 text-white shadow-sm overflow-hidden relative">
              <div className="absolute inset-0 bg-gradient-to-br from-blue-500 to-indigo-600" />
              <Cloud size={16} className="relative z-10" />
            </div>
            
            <AnimatePresence initial={false}>
              {!isCollapsed && (
                <motion.div 
                  initial={{ opacity: 0, width: 0 }}
                  animate={{ opacity: 1, width: "auto" }}
                  exit={{ opacity: 0, width: 0 }}
                  className="flex flex-col items-start min-w-0 overflow-hidden whitespace-nowrap"
                >
                  <span className="text-sm font-semibold truncate w-full text-left">CloudStoreX</span>
                  <span className="text-[11px] font-medium text-muted-foreground truncate w-full text-left flex items-center gap-1.5">
                    <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" /> Enterprise Plan
                  </span>
                </motion.div>
              )}
            </AnimatePresence>
          </div>
          
          <AnimatePresence initial={false}>
            {!isCollapsed && (
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
              >
                <ChevronsUpDown size={14} className="text-muted-foreground shrink-0" />
              </motion.div>
            )}
          </AnimatePresence>
      </DropdownMenuTrigger>
      
      <DropdownMenuContent align={isCollapsed ? "start" : "center"} side={isCollapsed ? "right" : "bottom"} className="w-64 rounded-xl" sideOffset={8}>
        <DropdownMenuLabel className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
          Active Workspace
        </DropdownMenuLabel>
        <DropdownMenuItem className="py-2.5 px-3 rounded-lg flex items-center justify-between bg-muted/30 focus:bg-muted/50 cursor-pointer">
          <div className="flex items-center gap-3">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-500 to-indigo-600 text-white shadow-sm">
              <Cloud size={16} />
            </div>
            <div className="flex flex-col">
              <span className="text-sm font-semibold text-foreground">CloudStoreX</span>
              <span className="text-[11px] text-muted-foreground">US East • 3 Providers</span>
            </div>
          </div>
          <Check size={16} className="text-primary" />
        </DropdownMenuItem>

        <DropdownMenuSeparator className="my-2" />

        <DropdownMenuItem className="py-2 px-3 rounded-lg flex items-center gap-3 cursor-pointer text-muted-foreground hover:text-foreground">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg border border-dashed bg-muted/30">
            <Plus size={16} />
          </div>
          <span className="text-sm font-medium">Create Workspace</span>
        </DropdownMenuItem>

        <DropdownMenuItem className="py-2 px-3 rounded-lg flex items-center gap-3 cursor-pointer text-muted-foreground hover:text-foreground">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-muted/30">
            <Settings size={16} />
          </div>
          <span className="text-sm font-medium">Workspace Settings</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
