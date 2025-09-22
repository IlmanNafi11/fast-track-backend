INSERT INTO users (email, password, name, is_active) 
VALUES ('user@example.com', '$2a$10$YourHashedPasswordHere', 'user example', true)
ON CONFLICT (email) DO NOTHING;