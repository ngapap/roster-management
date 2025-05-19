-- +goose Up
INSERT INTO users ( email, name, password, is_admin)
VALUES (
    'admin@roster.com',
    'SysAdmin',
    '$2a$10$YMgnaUDRDLMjIio/ObPPRuSofgRDSE8c23ubaofvI1ZVZrg3FOdL6', -- hashed 'Password.1'
    true
) ON CONFLICT (email) DO NOTHING;

-- +goose Down
DELETE FROM users WHERE email='admin@roster.com' AND name='SysAdmin';