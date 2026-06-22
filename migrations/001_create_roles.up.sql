CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(resource, action)
);

CREATE TABLE IF NOT EXISTS role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(role_id, permission_id)
);

-- Seed predefined roles
INSERT INTO roles (name, description) VALUES
    ('admin', 'Full access — can manage customers, orders, and staff'),
    ('staff', 'Limited access — can manage customers and orders')
ON CONFLICT (name) DO NOTHING;

-- Seed permissions
INSERT INTO permissions (resource, action, description) VALUES
    ('customers', 'read', 'View customers'),
    ('customers', 'create', 'Create new customer'),
    ('customers', 'update', 'Update customer details'),
    ('customers', 'delete', 'Delete customer'),
    ('orders', 'read', 'View orders'),
    ('orders', 'create', 'Create new order'),
    ('orders', 'update', 'Update order details'),
    ('orders', 'update_status', 'Change order status'),
    ('orders', 'delete', 'Delete order'),
    ('users', 'read', 'View staff users'),
    ('users', 'create', 'Create staff account'),
    ('users', 'update', 'Update staff account'),
    ('users', 'delete', 'Delete staff account'),
    ('reports', 'read', 'View reports and dashboard')
ON CONFLICT (resource, action) DO NOTHING;

-- Assign permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign permissions to staff role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'staff' AND p.resource IN ('customers', 'orders')
  AND p.action NOT IN ('delete', 'update_status')
ON CONFLICT (role_id, permission_id) DO NOTHING;

CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
CREATE INDEX IF NOT EXISTS idx_permissions_resource_action ON permissions(resource, action);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);
