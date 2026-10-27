-- 000002_schema.down.sql

DROP TABLE IF EXISTS rab_items;
DROP TABLE IF EXISTS volumes;
DROP TABLE IF EXISTS boq_items;
DROP TABLE IF EXISTS ahsp_items;
DROP TABLE IF EXISTS ahsps;
DROP TABLE IF EXISTS unit_prices;
DROP TABLE IF EXISTS schedules;
DROP TABLE IF EXISTS projects;
DELETE FROM users WHERE email = 'admin@musashi.local';

