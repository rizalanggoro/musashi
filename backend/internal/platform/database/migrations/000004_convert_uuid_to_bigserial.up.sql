-- 000004_convert_uuid_to_bigserial.up.sql
-- Transition database tables from UUID primary/foreign keys to BIGSERIAL auto-increment integers
-- and ensure deleted_at exists across all tables for GORM soft delete support.

-- 1. Drop existing foreign key constraints that reference UUID primary keys
ALTER TABLE IF EXISTS schedules DROP CONSTRAINT IF EXISTS schedules_project_id_fkey;
ALTER TABLE IF EXISTS ahsp_items DROP CONSTRAINT IF EXISTS ahsp_items_ahsp_id_fkey;
ALTER TABLE IF EXISTS ahsp_items DROP CONSTRAINT IF EXISTS ahsp_items_unit_price_id_fkey;
ALTER TABLE IF EXISTS boq_items DROP CONSTRAINT IF EXISTS boq_items_project_id_fkey;
ALTER TABLE IF EXISTS boq_items DROP CONSTRAINT IF EXISTS boq_items_ahsp_id_fkey;
ALTER TABLE IF EXISTS volumes DROP CONSTRAINT IF EXISTS volumes_boq_item_id_fkey;
ALTER TABLE IF EXISTS rab_items DROP CONSTRAINT IF EXISTS rab_items_project_id_fkey;
ALTER TABLE IF EXISTS rab_items DROP CONSTRAINT IF EXISTS rab_items_boq_item_id_fkey;

-- 2. USERS table: UUID -> BIGSERIAL + deleted_at
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE users DROP COLUMN IF EXISTS id;
ALTER TABLE users ADD COLUMN id bigserial PRIMARY KEY;
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at timestamptz;
CREATE INDEX IF NOT EXISTS users_deleted_at_idx ON users (deleted_at);
DROP INDEX IF EXISTS users_email_lower_idx;
CREATE UNIQUE INDEX users_email_lower_idx ON users (lower(email)) WHERE deleted_at IS NULL;

-- 3. PROJECTS table: UUID -> BIGSERIAL + deleted_at
ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_pkey;
ALTER TABLE projects DROP COLUMN IF EXISTS id;
ALTER TABLE projects ADD COLUMN id bigserial PRIMARY KEY;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS deleted_at timestamptz;
CREATE INDEX IF NOT EXISTS projects_deleted_at_idx ON projects (deleted_at);
DROP INDEX IF EXISTS projects_name_lower_idx;
CREATE UNIQUE INDEX projects_name_lower_idx ON projects (lower(name)) WHERE deleted_at IS NULL;

-- 4. SCHEDULES table: UUID -> BIGSERIAL + deleted_at
ALTER TABLE schedules DROP CONSTRAINT IF EXISTS schedules_pkey;
ALTER TABLE schedules DROP COLUMN IF EXISTS id;
ALTER TABLE schedules ADD COLUMN id bigserial PRIMARY KEY;
ALTER TABLE schedules DROP COLUMN IF EXISTS project_id;
ALTER TABLE schedules ADD COLUMN project_id bigint NOT NULL REFERENCES projects(id) ON DELETE CASCADE;
ALTER TABLE schedules ADD COLUMN IF NOT EXISTS deleted_at timestamptz;
CREATE INDEX IF NOT EXISTS schedules_project_id_idx ON schedules (project_id);
CREATE INDEX IF NOT EXISTS schedules_deleted_at_idx ON schedules (deleted_at);

-- 5. UNIT_PRICES table: UUID -> BIGSERIAL + deleted_at
ALTER TABLE unit_prices DROP CONSTRAINT IF EXISTS unit_prices_pkey;
ALTER TABLE unit_prices DROP COLUMN IF EXISTS id;
ALTER TABLE unit_prices ADD COLUMN id bigserial PRIMARY KEY;
ALTER TABLE unit_prices ADD COLUMN IF NOT EXISTS deleted_at timestamptz;
CREATE INDEX IF NOT EXISTS unit_prices_category_idx ON unit_prices (category);
CREATE INDEX IF NOT EXISTS unit_prices_deleted_at_idx ON unit_prices (deleted_at);
DROP INDEX IF EXISTS unit_prices_code_lower_idx;
CREATE UNIQUE INDEX unit_prices_code_lower_idx ON unit_prices (lower(code));

-- 6. AHSPS table: UUID -> BIGSERIAL + deleted_at
ALTER TABLE ahsps DROP CONSTRAINT IF EXISTS ahsps_pkey;
ALTER TABLE ahsps DROP COLUMN IF EXISTS id;
ALTER TABLE ahsps ADD COLUMN id bigserial PRIMARY KEY;
ALTER TABLE ahsps ADD COLUMN IF NOT EXISTS deleted_at timestamptz;
DROP INDEX IF EXISTS ahsps_code_lower_idx;
CREATE UNIQUE INDEX ahsps_code_lower_idx ON ahsps (lower(code));
CREATE INDEX IF NOT EXISTS ahsps_deleted_at_idx ON ahsps (deleted_at);

-- 7. AHSP_ITEMS table: UUID -> BIGSERIAL + deleted_at
ALTER TABLE ahsp_items DROP CONSTRAINT IF EXISTS ahsp_items_pkey;
ALTER TABLE ahsp_items DROP COLUMN IF EXISTS id;
ALTER TABLE ahsp_items ADD COLUMN id bigserial PRIMARY KEY;
ALTER TABLE ahsp_items DROP COLUMN IF EXISTS ahsp_id;
ALTER TABLE ahsp_items ADD COLUMN ahsp_id bigint NOT NULL REFERENCES ahsps(id) ON DELETE CASCADE;
ALTER TABLE ahsp_items DROP COLUMN IF EXISTS unit_price_id;
ALTER TABLE ahsp_items ADD COLUMN unit_price_id bigint REFERENCES unit_prices(id) ON DELETE RESTRICT;
ALTER TABLE ahsp_items ADD COLUMN IF NOT EXISTS deleted_at timestamptz;
CREATE INDEX IF NOT EXISTS ahsp_items_ahsp_id_idx ON ahsp_items (ahsp_id);
CREATE INDEX IF NOT EXISTS ahsp_items_unit_price_id_idx ON ahsp_items (unit_price_id);
CREATE INDEX IF NOT EXISTS ahsp_items_deleted_at_idx ON ahsp_items (deleted_at);

-- 8. BOQ_ITEMS table: UUID -> BIGSERIAL + deleted_at
ALTER TABLE boq_items DROP CONSTRAINT IF EXISTS boq_items_pkey;
ALTER TABLE boq_items DROP COLUMN IF EXISTS id;
ALTER TABLE boq_items ADD COLUMN id bigserial PRIMARY KEY;
ALTER TABLE boq_items DROP COLUMN IF EXISTS project_id;
ALTER TABLE boq_items ADD COLUMN project_id bigint NOT NULL REFERENCES projects(id) ON DELETE CASCADE;
ALTER TABLE boq_items DROP COLUMN IF EXISTS ahsp_id;
ALTER TABLE boq_items ADD COLUMN ahsp_id bigint REFERENCES ahsps(id) ON DELETE SET NULL;
ALTER TABLE boq_items ADD COLUMN IF NOT EXISTS deleted_at timestamptz;
CREATE INDEX IF NOT EXISTS boq_items_project_id_idx ON boq_items (project_id);
CREATE INDEX IF NOT EXISTS boq_items_deleted_at_idx ON boq_items (deleted_at);

-- 9. VOLUMES table: UUID -> BIGSERIAL + deleted_at
ALTER TABLE volumes DROP CONSTRAINT IF EXISTS volumes_pkey;
ALTER TABLE volumes DROP COLUMN IF EXISTS id;
ALTER TABLE volumes ADD COLUMN id bigserial PRIMARY KEY;
ALTER TABLE volumes DROP COLUMN IF EXISTS boq_item_id;
ALTER TABLE volumes ADD COLUMN boq_item_id bigint NOT NULL REFERENCES boq_items(id) ON DELETE CASCADE;
ALTER TABLE volumes ADD COLUMN IF NOT EXISTS deleted_at timestamptz;
CREATE INDEX IF NOT EXISTS volumes_boq_item_id_idx ON volumes (boq_item_id);
CREATE INDEX IF NOT EXISTS volumes_deleted_at_idx ON volumes (deleted_at);

-- 10. RAB_ITEMS table: UUID -> BIGSERIAL + deleted_at
ALTER TABLE rab_items DROP CONSTRAINT IF EXISTS rab_items_pkey;
ALTER TABLE rab_items DROP COLUMN IF EXISTS id;
ALTER TABLE rab_items ADD COLUMN id bigserial PRIMARY KEY;
ALTER TABLE rab_items DROP COLUMN IF EXISTS project_id;
ALTER TABLE rab_items ADD COLUMN project_id bigint NOT NULL REFERENCES projects(id) ON DELETE CASCADE;
ALTER TABLE rab_items DROP COLUMN IF EXISTS boq_item_id;
ALTER TABLE rab_items ADD COLUMN boq_item_id bigint REFERENCES boq_items(id) ON DELETE SET NULL;
ALTER TABLE rab_items ADD COLUMN IF NOT EXISTS deleted_at timestamptz;
CREATE INDEX IF NOT EXISTS rab_items_project_id_idx ON rab_items (project_id);
CREATE INDEX IF NOT EXISTS rab_items_deleted_at_idx ON rab_items (deleted_at);

-- Ensure default admin user is present with updated password hash
INSERT INTO users (email, name, password_hash, role)
VALUES (
    'admin@musashi.local',
    'Administrator',
    '$2a$10$i/3zZPOBK/Bw79ayDXks9u3YblN8XnEAQqruvWzi2uwzl4u/TQerS',
    'admin'
) ON CONFLICT (email) DO UPDATE SET password_hash = EXCLUDED.password_hash;

