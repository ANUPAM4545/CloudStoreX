# Role-Based Access Control (RBAC)

RBAC supports system-level and organization-level scoped roles via `RoleBindings`.

Inheritance is natively modeled where Organization Admins automatically receive escalated permissions across child Workspaces.
Wildcard matching is supported (e.g. `storage:object:*`).
