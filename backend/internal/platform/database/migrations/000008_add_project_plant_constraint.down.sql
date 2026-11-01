ALTER TABLE projects
    DROP CONSTRAINT IF EXISTS chk_projects_plant,
    ALTER COLUMN plant SET DEFAULT '';

