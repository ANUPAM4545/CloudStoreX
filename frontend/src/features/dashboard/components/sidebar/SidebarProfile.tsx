"use client";

import * as React from "react";
import { motion, AnimatePresence } from "framer-motion";
import { UserCircle, ChevronsUpDown, Settings, Key, Activity, LogOut } from "lucide-react";
import { useSidebarStore } from "@/store/sidebar";
import { useAuthStore } from "@/store/auth";
import { useRouter } from "next/navigation";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export function SidebarProfile() {
  const { isCollapsed } = useSidebarStore();
  const { user, logout } = useAuthStore();
  const router = useRouter();
  const [open, setOpen] = React.useState(false);

  const handleSignOut = () => {
    logout();
    router.push("/login");
  };

  return (
    <div className="p-3 mt-auto border-t border-border/40">
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
              <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-gradient-to-tr from-violet-500 to-fuchsia-500 text-white shadow-sm ring-1 ring-white/10 overflow-hidden relative">
                <span className="font-semibold text-xs relative z-10">
                  {user?.fullName?.charAt(0).toUpperCase() || 'U'}
                </span>
              </div>
              
              <AnimatePresence initial={false}>
                {!isCollapsed && (
                  <motion.div 
                    initial={{ opacity: 0, width: 0 }}
                    animate={{ opacity: 1, width: "auto" }}
                    exit={{ opacity: 0, width: 0 }}
                    className="flex flex-col items-start min-w-0 overflow-hidden whitespace-nowrap"
                  >
                    <span className="text-sm font-semibold truncate w-full text-left">{user?.fullName || 'User'}</span>
                    <span className="text-[11px] font-medium text-muted-foreground truncate w-full text-left">{user?.email || 'user@example.com'}</span>
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
        
        <DropdownMenuContent align={isCollapsed ? "start" : "center"} side={isCollapsed ? "right" : "bottom"} className="w-56 rounded-xl" sideOffset={8}>
          <DropdownMenuLabel className="font-normal py-2 px-3">
            <div className="flex flex-col space-y-1">
              <p className="text-sm font-medium leading-none">{user?.fullName || 'User'}</p>
              <p className="text-xs leading-none text-muted-foreground">{user?.email || 'user@example.com'}</p>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuGroup>
            <DropdownMenuItem className="py-2 cursor-pointer" onClick={() => router.push("/dashboard/settings/profile")}>
              <UserCircle className="mr-2 h-4 w-4 text-muted-foreground" />
              <span>Profile</span>
            </DropdownMenuItem>
            <DropdownMenuItem className="py-2 cursor-pointer" onClick={() => router.push("/dashboard/settings/configuration")}>
              <Settings className="mr-2 h-4 w-4 text-muted-foreground" />
              <span>Settings</span>
            </DropdownMenuItem>
            <DropdownMenuItem className="py-2 cursor-pointer" onClick={() => router.push("/dashboard/settings/api-keys")}>
              <Key className="mr-2 h-4 w-4 text-muted-foreground" />
              <span>API Keys</span>
            </DropdownMenuItem>
            <DropdownMenuItem className="py-2 cursor-pointer" onClick={() => router.push("/dashboard/settings/sessions")}>
              <Activity className="mr-2 h-4 w-4 text-muted-foreground" />
              <span>Sessions</span>
            </DropdownMenuItem>
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <DropdownMenuItem className="py-2 text-rose-500 focus:text-rose-600 focus:bg-rose-500/10 cursor-pointer" onClick={handleSignOut}>
            <LogOut className="mr-2 h-4 w-4" />
            <span>Sign out</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}
