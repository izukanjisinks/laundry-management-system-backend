-- Seed admin user (password: Admin1234! — change after first login)
INSERT INTO users (full_name, email, password, role_id, is_active)
SELECT
    'System Admin',
    'admin@washpoint.app',
    '$2a$12$1GVAho02IXJ.1Jrof5OSpeDTQkTpQCgs0XtmH3Aj8voVvPSvamFim',
    r.id,
    true
FROM roles r
WHERE r.name = 'admin'
ON CONFLICT (email) DO NOTHING;
