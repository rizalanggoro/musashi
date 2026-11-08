-- 000011_add_formula_to_volumes.down.sql
ALTER TABLE volumes DROP COLUMN IF EXISTS formula;

