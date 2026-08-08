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
  Image as ImageIcon,
  Clock,
  Users,
  Star,
  Trash2,
  TerminalSquare,
  Key,
  Briefcase
} from "lucide-react";
import { ElementType } from "react";

export type ExperienceType = 'personal' | 'developer' | 'enterprise';

export interface NavItem {
  id: string;
  title: string;
  description?: string;
  href: string;
  icon: ElementType;
  group?: string;
  order: number;
  badge?: string;
  featureFlag?: string;
  requiredExperience?: ExperienceType[];
  requiredRole?: string[];
  requiredPermission?: string[];
  comingSoon?: boolean;
  children?: NavItem[];
}

export const NAVIGATION_CONFIG: NavItem[] = [
  // ==========================================
  // PERSONAL EXPERIENCE
  // ==========================================
  { id: "p-files", title: "Files", href: "/personal/files", icon: FileText, group: "Storage", order: 1, requiredExperience: ["personal"] },
  { id: "p-photos", title: "Photos", href: "/personal/photos", icon: ImageIcon, group: "Storage", order: 2, requiredExperience: ["personal"] },
  { id: "p-recent", title: "Recent", href: "/personal/recent", icon: Clock, group: "Storage", order: 3, requiredExperience: ["personal"] },
  { id: "p-shared", title: "Shared", href: "/personal/shared", icon: Users, group: "Storage", order: 4, requiredExperience: ["personal"] },
  { id: "p-favorites", title: "Favorites", href: "/personal/favorites", icon: Star, group: "Storage", order: 5, requiredExperience: ["personal"] },
  { id: "p-ai-search", title: "AI Search", href: "/personal/search", icon: Sparkles, group: "Intelligence", order: 6, requiredExperience: ["personal"], badge: "New" },
  { id: "p-trash", title: "Trash", href: "/personal/trash", icon: Trash2, group: "System", order: 7, requiredExperience: ["personal"] },
  { id: "p-settings", title: "Settings", href: "/personal/settings", icon: Settings2, group: "System", order: 8, requiredExperience: ["personal"] },

  // ==========================================
  // DEVELOPER EXPERIENCE
  // ==========================================
  { id: "d-overview", title: "Overview", href: "/developer/overview", icon: LayoutDashboard, group: "Platform", order: 1, requiredExperience: ["developer"] },
  { id: "d-buckets", title: "Buckets", href: "/developer/buckets", icon: Database, group: "Platform", order: 2, requiredExperience: ["developer"] },
  { id: "d-objects", title: "Objects", href: "/developer/objects", icon: FileText, group: "Platform", order: 3, requiredExperience: ["developer"] },
  { id: "d-providers", title: "Providers", href: "/developer/providers", icon: Cloud, group: "Platform", order: 4, requiredExperience: ["developer"] },
  { id: "d-policies", title: "Policies", href: "/developer/policies", icon: GitBranch, group: "Platform", order: 5, requiredExperience: ["developer"] },
  { id: "d-apikeys", title: "API Keys", href: "/developer/keys", icon: Key, group: "Developers", order: 6, requiredExperience: ["developer"] },
  { id: "d-sdks", title: "SDKs", href: "/developer/sdks", icon: Briefcase, group: "Developers", order: 7, requiredExperience: ["developer"] },
  { id: "d-cli", title: "CLI", href: "/developer/cli", icon: TerminalSquare, group: "Developers", order: 8, requiredExperience: ["developer"] },
  { id: "d-monitoring", title: "Monitoring", href: "/developer/monitoring", icon: Activity, group: "Observability", order: 9, requiredExperience: ["developer"] },
  { id: "d-analytics", title: "Analytics", href: "/developer/analytics", icon: ChartColumn, group: "Observability", order: 10, requiredExperience: ["developer"] },
  { id: "d-ai", title: "AI Assistant", href: "/developer/ai", icon: Sparkles, group: "Intelligence", order: 11, requiredExperience: ["developer"] },

  // ==========================================
  // ENTERPRISE EXPERIENCE
  // ==========================================
  { id: "e-overview", title: "Overview", href: "/enterprise/overview", icon: LayoutDashboard, group: "Platform", order: 1, requiredExperience: ["enterprise"] },
  { id: "e-orgs", title: "Organizations", href: "/enterprise/organizations", icon: Building, group: "Access Management", order: 2, requiredExperience: ["enterprise"] },
  { id: "e-teams", title: "Teams", href: "/enterprise/teams", icon: Users, group: "Access Management", order: 3, requiredExperience: ["enterprise"] },
  { id: "e-users", title: "Users", href: "/enterprise/users", icon: UserCircle, group: "Access Management", order: 4, requiredExperience: ["enterprise"] },
  { id: "e-rbac", title: "RBAC", href: "/enterprise/rbac", icon: Shield, group: "Access Management", order: 5, requiredExperience: ["enterprise"] },
  
  { id: "e-providers", title: "Providers", href: "/enterprise/providers", icon: Cloud, group: "Infrastructure", order: 6, requiredExperience: ["enterprise"] },
  { id: "e-policies", title: "Policies", href: "/enterprise/policies", icon: GitBranch, group: "Infrastructure", order: 7, requiredExperience: ["enterprise"] },
  { id: "e-replication", title: "Replication", href: "/enterprise/replication", icon: Copy, group: "Infrastructure", order: 8, requiredExperience: ["enterprise"] },
  { id: "e-dr", title: "Disaster Recovery", href: "/enterprise/dr", icon: LifeBuoy, group: "Infrastructure", order: 9, requiredExperience: ["enterprise"] },
  
  { id: "e-security", title: "Security", href: "/enterprise/security", icon: Shield, group: "Governance", order: 10, requiredExperience: ["enterprise"] },
  { id: "e-audit", title: "Audit", href: "/enterprise/audit", icon: FileSearch, group: "Governance", order: 11, requiredExperience: ["enterprise"] },
  
  { id: "e-analytics", title: "Analytics", href: "/enterprise/analytics", icon: ChartColumn, group: "Intelligence", order: 12, requiredExperience: ["enterprise"] },
  { id: "e-ai", title: "AI Center", href: "/enterprise/ai", icon: Sparkles, group: "Intelligence", order: 13, requiredExperience: ["enterprise"] },
  
  { id: "e-billing", title: "Billing", href: "/enterprise/billing", icon: Database, group: "Administration", order: 14, requiredExperience: ["enterprise"] },
  { id: "e-settings", title: "Settings", href: "/enterprise/settings/profile", icon: Settings2, group: "Administration", order: 15, requiredExperience: ["enterprise"] },
];

/**
 * Helper to get grouped navigation for a specific experience and role
 */
export function getNavigation(
  experience: ExperienceType,
  userRoles: string[] = [],
  userPermissions: string[] = []
) {
  // Filter by experience
  let filtered = NAVIGATION_CONFIG.filter((item) => {
    if (item.requiredExperience && !item.requiredExperience.includes(experience)) {
      return false;
    }
    // Future RBAC filtering can be added here using userRoles and userPermissions
    return true;
  });

  // Sort by order
  filtered.sort((a, b) => a.order - b.order);

  // Group items
  const groups: Record<string, NavItem[]> = {};
  
  filtered.forEach(item => {
    const groupName = item.group || "General";
    if (!groups[groupName]) {
      groups[groupName] = [];
    }
    groups[groupName].push(item);
  });

  // Convert to array of groups
  return Object.keys(groups).map(key => ({
    label: key,
    items: groups[key]
  }));
}
