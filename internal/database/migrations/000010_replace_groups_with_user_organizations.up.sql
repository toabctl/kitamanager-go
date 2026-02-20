-- Create user_organizations table with direct user-org-role relationship
CREATE TABLE user_organizations (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),
    PRIMARY KEY (user_id, organization_id)
);

-- Migrate data: for each user+org pair, take the highest-precedence role
INSERT INTO user_organizations (user_id, organization_id, role, created_at, created_by)
SELECT DISTINCT ON (ug.user_id, g.organization_id)
    ug.user_id,
    g.organization_id,
    ug.role,
    ug.created_at,
    ug.created_by
FROM user_groups ug
JOIN groups g ON g.id = ug.group_id
ORDER BY ug.user_id, g.organization_id,
    CASE ug.role
        WHEN 'admin' THEN 3
        WHEN 'manager' THEN 2
        WHEN 'member' THEN 1
        ELSE 0
    END DESC,
    ug.created_at ASC;

-- Drop old tables
DROP TABLE IF EXISTS user_groups;
DROP TABLE IF EXISTS groups;
