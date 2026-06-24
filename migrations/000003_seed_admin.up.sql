-- Seed admin user (idempotent).
-- Password: Pretest@2025 (bcrypt cost 12)
INSERT INTO users (id, username, password_hash, role)
SELECT
    '00000000-0000-0000-0000-000000000001',
    'admin',
    '$2a$12$ilFIA9G0SnjRt16fPPAaJ.BBg.SkHJ3JXJFQyrITH7d/jR9jbgEtu',
    'ADMIN'
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE username = 'admin'
);
