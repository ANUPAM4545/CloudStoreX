import { create } from "zustand";
import { persist } from "zustand/middleware";
import { ExperienceType } from "@/config/navigation";

export interface Workspace {
  id: string;
  name: string;
  experience: ExperienceType;
  plan: string;
  region: string;
}

interface WorkspaceState {
  activeWorkspace: Workspace | null;
  workspaces: Workspace[];
  setActiveWorkspace: (workspaceId: string) => void;
  addWorkspace: (workspace: Workspace) => void;
  updateExperience: (experience: ExperienceType) => void;
}

// Mock initial data
const INITIAL_WORKSPACES: Workspace[] = [
  {
    id: "ws-personal",
    name: "Personal Cloud",
    experience: "personal",
    plan: "Free Tier",
    region: "Global",
  },
  {
    id: "ws-dev",
    name: "Hackathon Project",
    experience: "developer",
    plan: "Pro Plan",
    region: "US East",
  },
  {
    id: "ws-ent",
    name: "CloudStoreX Corp",
    experience: "enterprise",
    plan: "Enterprise",
    region: "US East (Multi-Zone)",
  },
];

export const useWorkspaceStore = create<WorkspaceState>()(
  persist(
    (set) => ({
      // Start with the enterprise workspace as default to match existing flow,
      // but in reality this would be fetched from the backend.
      activeWorkspace: INITIAL_WORKSPACES[2],
      workspaces: INITIAL_WORKSPACES,

      setActiveWorkspace: (workspaceId) =>
        set((state) => {
          const ws = state.workspaces.find((w) => w.id === workspaceId);
          if (ws) {
            return { activeWorkspace: ws };
          }
          return state;
        }),

      addWorkspace: (workspace) =>
        set((state) => ({
          workspaces: [...state.workspaces, workspace],
          activeWorkspace: workspace, // Auto-switch to new workspace
        })),

      updateExperience: (experience) =>
        set((state) => {
          if (!state.activeWorkspace) return state;
          
          const updated = { ...state.activeWorkspace, experience };
          const updatedWorkspaces = state.workspaces.map((w) => 
            w.id === updated.id ? updated : w
          );
          
          return {
            activeWorkspace: updated,
            workspaces: updatedWorkspaces,
          };
        }),
    }),
    {
      name: "cloudstorex-workspace-state",
    }
  )
);
