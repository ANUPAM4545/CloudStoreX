import { create } from "zustand";
import { persist } from "zustand/middleware";

interface SidebarState {
  isCollapsed: boolean;
  toggleCollapse: () => void;
  setCollapsed: (collapsed: boolean) => void;
  
  // Mobile drawer state
  isMobileOpen: boolean;
  setMobileOpen: (open: boolean) => void;
  
  // Command palette state (placeholder for now)
  isCommandPaletteOpen: boolean;
  setCommandPaletteOpen: (open: boolean) => void;
}

export const useSidebarStore = create<SidebarState>()(
  persist(
    (set) => ({
      isCollapsed: false,
      toggleCollapse: () => set((state) => ({ isCollapsed: !state.isCollapsed })),
      setCollapsed: (collapsed) => set({ isCollapsed: collapsed }),
      
      isMobileOpen: false,
      setMobileOpen: (open) => set({ isMobileOpen: open }),
      
      isCommandPaletteOpen: false,
      setCommandPaletteOpen: (open) => set({ isCommandPaletteOpen: open }),
    }),
    {
      name: "cloudstorex-sidebar-state", // unique name for localStorage
      partialize: (state) => ({ isCollapsed: state.isCollapsed }), // Only persist isCollapsed
    }
  )
);
