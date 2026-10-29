-- 000002_schema.up.sql

CREATE TABLE projects (
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name                text NOT NULL,
    section             text NOT NULL DEFAULT '',
    plant               text NOT NULL DEFAULT '',
    implementation_plan date,
    budget              bigint NOT NULL DEFAULT 0,
    information         text NOT NULL DEFAULT '',
    status              text NOT NULL DEFAULT 'planning' CHECK (status IN ('planning', 'on_going', 'completed', 'cancelled')),
    progress_spb        text NOT NULL DEFAULT 'pending' CHECK (progress_spb IN ('pending', 'done')),
    progress_rab        text NOT NULL DEFAULT 'pending' CHECK (progress_rab IN ('pending', 'done')),
    progress_xxx        text NOT NULL DEFAULT 'pending' CHECK (progress_xxx IN ('pending', 'done')),
    deleted_at          timestamptz,
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX projects_name_lower_idx ON projects (lower(name)) WHERE deleted_at IS NULL;
CREATE INDEX projects_deleted_at_idx ON projects (deleted_at);

CREATE TABLE schedules (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    code        text NOT NULL DEFAULT '',
    work_item   text NOT NULL,
    start_date  date NOT NULL,
    end_date    date NOT NULL,
    duration    int NOT NULL DEFAULT 1,
    weight      numeric(6,2) NOT NULL DEFAULT 0,
    progress    numeric(6,2) NOT NULL DEFAULT 0,
    status      text NOT NULL DEFAULT 'planning' CHECK (status IN ('planning', 'on_going', 'done', 'late')),
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX schedules_project_id_idx ON schedules (project_id);

CREATE TABLE unit_prices (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code        text NOT NULL UNIQUE,
    category    text NOT NULL CHECK (category IN ('LABOR', 'MATERIALS', 'EQUIPMENT', 'OTHER')),
    work_item   text NOT NULL,
    unit        text NOT NULL,
    unit_price  bigint NOT NULL DEFAULT 0,
    information text NOT NULL DEFAULT '',
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX unit_prices_category_idx ON unit_prices (category);
CREATE UNIQUE INDEX unit_prices_code_lower_idx ON unit_prices (lower(code));

CREATE TABLE ahsps (
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code                text NOT NULL UNIQUE,
    work_item           text NOT NULL,
    unit                text NOT NULL,
    overhead_profit_pct numeric(6,2) NOT NULL DEFAULT 10.00,
    total_labor         bigint NOT NULL DEFAULT 0,
    total_material      bigint NOT NULL DEFAULT 0,
    total_equipment     bigint NOT NULL DEFAULT 0,
    total_cost          bigint NOT NULL DEFAULT 0,
    unit_price_of_work  bigint NOT NULL DEFAULT 0,
    status              text NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'final')),
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ahsps_code_lower_idx ON ahsps (lower(code));

CREATE TABLE ahsp_items (
    id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    ahsp_id        uuid NOT NULL REFERENCES ahsps(id) ON DELETE CASCADE,
    category       text NOT NULL CHECK (category IN ('LABOR', 'MATERIALS', 'EQUIPMENT')),
    sort           int NOT NULL DEFAULT 0,
    description    text NOT NULL DEFAULT '',
    code           text NOT NULL DEFAULT '',
    unit           text NOT NULL DEFAULT '',
    coefficient    numeric(14,4) NOT NULL DEFAULT 0,
    unit_price     bigint NOT NULL DEFAULT 0,
    cost           bigint NOT NULL DEFAULT 0,
    unit_price_id  uuid REFERENCES unit_prices(id) ON DELETE RESTRICT,
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX ahsp_items_ahsp_id_idx ON ahsp_items (ahsp_id);
CREATE INDEX ahsp_items_unit_price_id_idx ON ahsp_items (unit_price_id);

CREATE TABLE boq_items (
    id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id     uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    ahsp_id        uuid REFERENCES ahsps(id) ON DELETE SET NULL,
    work_item      text NOT NULL,
    unit           text NOT NULL,
    planned_volume numeric(14,4) NOT NULL DEFAULT 0,
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX boq_items_project_id_idx ON boq_items (project_id);

CREATE TABLE volumes (
    id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    boq_item_id    uuid NOT NULL REFERENCES boq_items(id) ON DELETE CASCADE,
    date           date,
    length         numeric(14,4),
    width          numeric(14,4),
    height         numeric(14,4),
    quantity       numeric(14,4) NOT NULL DEFAULT 0,
    unit_weight    numeric(14,4) NOT NULL DEFAULT 1,
    total_quantity numeric(14,4) NOT NULL DEFAULT 0,
    description    text NOT NULL DEFAULT '',
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX volumes_boq_item_id_idx ON volumes (boq_item_id);

CREATE TABLE rab_items (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    boq_item_id uuid REFERENCES boq_items(id) ON DELETE SET NULL,
    work_item   text NOT NULL,
    unit        text NOT NULL,
    volume      numeric(14,4) NOT NULL DEFAULT 0,
    unit_price  bigint NOT NULL DEFAULT 0,
    total_cost  bigint NOT NULL DEFAULT 0,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX rab_items_project_id_idx ON rab_items (project_id);

-- Default admin user: admin@musashi.local / admin123
INSERT INTO users (id, email, name, password_hash, role)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'admin@musashi.local',
    'Administrator',
    '$2a$10$i/3zZPOBK/Bw79ayDXks9u3YblN8XnEAQqruvWzi2uwzl4u/TQerS',
    'admin'
) ON CONFLICT (email) DO UPDATE SET password_hash = EXCLUDED.password_hash;
