-- 000011_add_formula_to_volumes.up.sql
ALTER TABLE volumes ADD COLUMN IF NOT EXISTS formula text NOT NULL DEFAULT '';

