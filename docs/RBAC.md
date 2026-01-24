# Role-Based Access Control (RBAC)

KitaManager uses [Casbin](https://casbin.org/) for role-based access control with organization-level multi-tenancy.

## Overview

The RBAC system allows:
- **Organization-scoped roles**: Users can have different roles in different organizations
- **Hierarchical permissions**: Roles define what actions users can perform on resources
- **Superadmin override**: Superadmins have full access across all organizations

## Roles

| Role | Description |
|------|-------------|
| `superadmin` | Full system access. Can create/delete organizations and access all resources across all organizations. |
| `admin` | Full access within their assigned organization(s). Cannot create/delete organizations. |
| `manager` | Operational access. Can manage employees, children, and contracts. Read-only access to users and groups. |

## Resources

| Resource | Description |
|----------|-------------|
| `organizations` | Kita/daycare centers |
| `employees` | Staff members |
| `children` | Enrolled children |
| `employee_contracts` | Employment contracts |
| `child_contracts` | Enrollment contracts |
| `users` | System users |
| `groups` | User groups |

## Actions

| Action | Description |
|--------|-------------|
| `create` | Create new records |
| `read` | View records |
| `update` | Modify existing records |
| `delete` | Remove records |

## Permission Matrix

### Superadmin

| Resource | Create | Read | Update | Delete |
|----------|--------|------|--------|--------|
| organizations | ✅ | ✅ | ✅ | ✅ |
| employees | ✅ | ✅ | ✅ | ✅ |
| children | ✅ | ✅ | ✅ | ✅ |
| employee_contracts | ✅ | ✅ | ✅ | ✅ |
| child_contracts | ✅ | ✅ | ✅ | ✅ |
| users | ✅ | ✅ | ✅ | ✅ |
| groups | ✅ | ✅ | ✅ | ✅ |

### Admin (within assigned organization)

| Resource | Create | Read | Update | Delete |
|----------|--------|------|--------|--------|
| organizations | ❌ | ✅ | ✅ | ❌ |
| employees | ✅ | ✅ | ✅ | ✅ |
| children | ✅ | ✅ | ✅ | ✅ |
| employee_contracts | ✅ | ✅ | ✅ | ✅ |
| child_contracts | ✅ | ✅ | ✅ | ✅ |
| users | ✅ | ✅ | ✅ | ✅ |
| groups | ✅ | ✅ | ✅ | ✅ |

### Manager (within assigned organization)

| Resource | Create | Read | Update | Delete |
|----------|--------|------|--------|--------|
| organizations | ❌ | ✅ | ❌ | ❌ |
| employees | ✅ | ✅ | ✅ | ✅ |
| children | ✅ | ✅ | ✅ | ✅ |
| employee_contracts | ✅ | ✅ | ✅ | ✅ |
| child_contracts | ✅ | ✅ | ✅ | ✅ |
| users | ❌ | ✅ | ❌ | ❌ |
| groups | ❌ | ✅ | ❌ | ❌ |

## API Endpoints

### Organization-Scoped Endpoints

Resources that belong to an organization use the URL pattern:
```
/api/v1/organizations/{orgId}/employees
/api/v1/organizations/{orgId}/children
```

These endpoints check if the user has the required permission for the specified organization.

### Global Endpoints

Some resources (users, groups) are global and not organization-scoped:
```
/api/v1/users
/api/v1/groups
```

### Superadmin-Only Endpoints

```
POST   /api/v1/organizations      # Create organization
DELETE /api/v1/organizations/{id} # Delete organization
```

## Casbin Configuration

### Model (`configs/rbac_model.conf`)

```ini
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && (p.dom == r.dom || p.dom == "*") && r.obj == p.obj && r.act == p.act
```

- `sub`: Subject (user ID, e.g., `user:1`)
- `dom`: Domain (organization ID, e.g., `org:1` or `*` for all)
- `obj`: Object (resource name)
- `act`: Action (create, read, update, delete)

### Policy Storage

Policies are stored in the PostgreSQL database using Casbin's GORM adapter. The table `casbin_rule` is automatically created.

## Usage

### Assigning Roles

```go
// Assign admin role to user 1 in organization 1
enforcer.AssignRole(1, "admin", 1)

// Assign superadmin (global access)
enforcer.AssignSuperAdmin(1)
```

### Checking Permissions

```go
// Check if user 1 can create employees in organization 1
allowed, _ := enforcer.CheckPermission(1, 1, "employees", "create")
```

### Initial Setup

On first run, seed the default policies:

```bash
SEED_RBAC_POLICIES=true ./kitamanager-api
```

This creates the role-permission mappings. You still need to assign roles to users.

### Creating the First Superadmin

After seeding policies, create a user and assign superadmin role via database or API:

```sql
-- Assuming user ID 1 exists
INSERT INTO casbin_rule (ptype, v0, v1, v2)
VALUES ('g', 'user:1', 'superadmin', '*');
```

Or programmatically:

```go
enforcer.AssignSuperAdmin(userID)
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `RBAC_MODEL_PATH` | Path to Casbin model config | `configs/rbac_model.conf` |
| `SEED_RBAC_POLICIES` | Seed default policies on startup | `false` |

## Testing

Run RBAC tests:

```bash
go test ./internal/rbac/... -v
go test ./internal/middleware/... -v -run Authorization
```
