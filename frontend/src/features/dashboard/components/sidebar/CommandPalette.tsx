"use client";

import * as React from "react";
import { useSidebarStore } from "@/store/sidebar";
import {
  Dialog,
  DialogContent,
} from "@/components/ui/dialog";
import { Search, FolderPlus, Cloud, Users, Sparkles, FileText, Activity } from "lucide-react";

export function CommandPalette() {
  const { isCommandPaletteOpen, setCommandPaletteOpen } = useSidebarStore();

  // Keyboard shortcut listener
  React.useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setCommandPaletteOpen(!isCommandPaletteOpen);
      }
    };

    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, [isCommandPaletteOpen, setCommandPaletteOpen]);

  return (
    <Dialog open={isCommandPaletteOpen} onOpenChange={setCommandPaletteOpen}>
      <DialogContent className="overflow-hidden p-0 rounded-2xl shadow-2xl border-border/40 sm:max-w-[600px]">
        
        <div className="flex items-center border-b px-4 py-3">
          <Search className="mr-3 h-5 w-5 text-muted-foreground" />
          <input 
            className="flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50"
            placeholder="Search resources, commands, or settings..." 
            autoFocus
          />
          <kbd className="ml-3 hidden sm:inline-flex h-6 select-none items-center gap-1 rounded border bg-muted px-2 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
            ESC
          </kbd>
        </div>

        <div className="max-h-[300px] overflow-y-auto p-2 scrollbar-none">
          <div className="px-2 py-1.5 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
            Quick Actions
          </div>
          
          <button className="relative flex cursor-pointer select-none items-center rounded-lg px-3 py-2.5 text-sm outline-none hover:bg-muted w-full transition-colors group">
            <FolderPlus className="mr-3 h-4 w-4 text-primary" />
            <span className="font-medium">Create new bucket</span>
          </button>
          
          <button className="relative flex cursor-pointer select-none items-center rounded-lg px-3 py-2.5 text-sm outline-none hover:bg-muted w-full transition-colors group">
            <Cloud className="mr-3 h-4 w-4 text-emerald-500" />
            <span className="font-medium">Connect storage provider</span>
          </button>
          
          <button className="relative flex cursor-pointer select-none items-center rounded-lg px-3 py-2.5 text-sm outline-none hover:bg-muted w-full transition-colors group">
            <Users className="mr-3 h-4 w-4 text-amber-500" />
            <span className="font-medium">Invite team member</span>
          </button>

          <div className="px-2 py-1.5 mt-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
            Navigation
          </div>
          
          <button className="relative flex cursor-pointer select-none items-center rounded-lg px-3 py-2.5 text-sm outline-none hover:bg-muted w-full transition-colors group text-muted-foreground hover:text-foreground">
            <Sparkles className="mr-3 h-4 w-4" />
            <span>Open AI Intelligence Center</span>
          </button>
          
          <button className="relative flex cursor-pointer select-none items-center rounded-lg px-3 py-2.5 text-sm outline-none hover:bg-muted w-full transition-colors group text-muted-foreground hover:text-foreground">
            <Activity className="mr-3 h-4 w-4" />
            <span>View Monitoring Dashboard</span>
          </button>
          
          <button className="relative flex cursor-pointer select-none items-center rounded-lg px-3 py-2.5 text-sm outline-none hover:bg-muted w-full transition-colors group text-muted-foreground hover:text-foreground">
            <FileText className="mr-3 h-4 w-4" />
            <span>Browse Object Explorer</span>
          </button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
