-- 000003_fix_admin_password.up.sql
-- Update default admin password hash to match "admin123"
UPDATE users
SET password_hash = '$2a$10$i/3zZPOBK/Bw79ayDXks9u3YblN8XnEAQqruvWzi2uwzl4u/TQerS'
WHERE email = 'admin@musashi.local';
