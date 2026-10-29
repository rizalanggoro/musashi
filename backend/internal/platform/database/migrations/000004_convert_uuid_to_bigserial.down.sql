-- 000004_convert_uuid_to_bigserial.down.sql
-- Revert tables back to UUID primary/foreign keys

ALTER TABLE IF EXISTS schedules DROP CONSTRAINT IF EXISTS schedules_project_id_fkey;
ALTER TABLE IF EXISTS ahsp_items DROP CONSTRAINT IF EXISTS ahsp_items_ahsp_id_fkey;
ALTER TABLE IF EXISTS ahsp_items DROP CONSTRAINT IF EXISTS ahsp_items_unit_price_id_fkey;
ALTER TABLE IF EXISTS boq_items DROP CONSTRAINT IF EXISTS boq_items_project_id_fkey;
ALTER TABLE IF EXISTS boq_items DROP CONSTRAINT IF EXISTS boq_items_ahsp_id_fkey;
ALTER TABLE IF EXISTS volumes DROP CONSTRAINT IF EXISTS volumes_boq_item_id_fkey;
ALTER TABLE IF EXISTS rab_items DROP CONSTRAINT IF EXISTS rab_items_project_id_fkey;
ALTER TABLE IF EXISTS rab_items DROP CONSTRAINT IF EXISTS rab_items_boq_item_id_fkey;

-- Revert users
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE users DROP COLUMN IF EXISTS id;
ALTER TABLE users ADD COLUMN id uuid PRIMARY KEY DEFAULT gen_random_uuid();
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;
DROP INDEX IF EXISTS users_deleted_at_idx;
DROP INDEX IF EXISTS users_email_lower_idx;
CREATE UNIQUE INDEX users_email_lower_idx ON users (lower(email));

-- Revert projects
ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_pkey;
ALTER TABLE projects DROP COLUMN IF EXISTS id;
ALTER TABLE projects ADD COLUMN id uuid PRIMARY KEY DEFAULT gen_random_uuid();

-- Revert schedules
ALTER TABLE schedules DROP CONSTRAINT IF EXISTS schedules_pkey;
ALTER TABLE schedules DROP COLUMN IF EXISTS id;
ALTER TABLE schedules ADD COLUMN id uuid PRIMARY KEY DEFAULT gen_random_uuid();
ALTER TABLE schedules DROP COLUMN IF EXISTS project_id;
ALTER TABLE schedules ADD COLUMN project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE;
ALTER TABLE schedules DROP COLUMN IF EXISTS deleted_at;
DROP INDEX IF EXISTS schedules_deleted_at_idx;

-- Revert unit_prices
ALTER TABLE unit_prices DROP CONSTRAINT IF EXISTS unit_prices_pkey;
ALTER TABLE unit_prices DROP COLUMN IF EXISTS id;
ALTER TABLE unit_prices ADD COLUMN id uuid PRIMARY KEY DEFAULT gen_random_uuid();
ALTER TABLE unit_prices DROP COLUMN IF EXISTS deleted_at;
DROP INDEX IF EXISTS unit_prices_deleted_at_idx;

-- Revert ahsps
ALTER TABLE ahsps DROP CONSTRAINT IF EXISTS ahsps_pkey;
ALTER TABLE ahsps DROP COLUMN IF EXISTS id;
ALTER TABLE ahsps ADD COLUMN id uuid PRIMARY KEY DEFAULT gen_random_uuid();
ALTER TABLE ahsps DROP COLUMN IF EXISTS deleted_at;
DROP INDEX IF EXISTS ahsps_deleted_at_idx;

-- Revert ahsp_items
ALTER TABLE ahsp_items DROP CONSTRAINT IF EXISTS ahsp_items_pkey;
ALTER TABLE ahsp_items DROP COLUMN IF EXISTS id;
ALTER TABLE ahsp_items ADD COLUMN id uuid PRIMARY KEY DEFAULT gen_random_uuid();
ALTER TABLE ahsp_items DROP COLUMN IF EXISTS ahsp_id;
ALTER TABLE ahsp_items ADD COLUMN ahsp_id uuid NOT NULL REFERENCES ahsps(id) ON DELETE CASCADE;
ALTER TABLE ahsp_items DROP COLUMN IF EXISTS unit_price_id;
ALTER TABLE ahsp_items ADD COLUMN unit_price_id uuid REFERENCES unit_prices(id) ON DELETE RESTRICT;
ALTER TABLE ahsp_items DROP COLUMN IF EXISTS deleted_at;
DROP INDEX IF EXISTS ahsp_items_deleted_at_idx;

-- Revert boq_items
ALTER TABLE boq_items DROP CONSTRAINT IF EXISTS boq_items_pkey;
ALTER TABLE boq_items DROP COLUMN IF EXISTS id;
ALTER TABLE boq_items ADD COLUMN id uuid PRIMARY KEY DEFAULT gen_random_uuid();
ALTER TABLE boq_items DROP COLUMN IF EXISTS project_id;
ALTER TABLE boq_items ADD COLUMN project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE;
ALTER TABLE boq_items DROP COLUMN IF EXISTS ahsp_id;
ALTER TABLE boq_items ADD COLUMN ahsp_id uuid REFERENCES ahsps(id) ON DELETE SET NULL;
ALTER TABLE boq_items DROP COLUMN IF EXISTS deleted_at;
DROP INDEX IF EXISTS boq_items_deleted_at_idx;

-- Revert volumes
ALTER TABLE volumes DROP CONSTRAINT IF EXISTS volumes_pkey;
ALTER TABLE volumes DROP COLUMN IF EXISTS id;
ALTER TABLE volumes ADD COLUMN id uuid PRIMARY KEY DEFAULT gen_random_uuid();
ALTER TABLE volumes DROP COLUMN IF EXISTS boq_item_id;
ALTER TABLE volumes ADD COLUMN boq_item_id uuid NOT NULL REFERENCES boq_items(id) ON DELETE CASCADE;
ALTER TABLE volumes DROP COLUMN IF EXISTS deleted_at;
DROP INDEX IF EXISTS volumes_deleted_at_idx;

-- Revert rab_items
ALTER TABLE rab_items DROP CONSTRAINT IF EXISTS rab_items_pkey;
ALTER TABLE rab_items DROP COLUMN IF EXISTS id;
ALTER TABLE rab_items ADD COLUMN id uuid PRIMARY KEY DEFAULT gen_random_uuid();
ALTER TABLE rab_items DROP COLUMN IF EXISTS project_id;
ALTER TABLE rab_items ADD COLUMN project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE;
ALTER TABLE rab_items DROP COLUMN IF EXISTS boq_item_id;
ALTER TABLE rab_items ADD COLUMN boq_item_id uuid REFERENCES boq_items(id) ON DELETE SET NULL;
ALTER TABLE rab_items DROP COLUMN IF EXISTS deleted_at;
DROP INDEX IF EXISTS rab_items_deleted_at_idx;

