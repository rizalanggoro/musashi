ALTER TABLE projects
    DROP CONSTRAINT IF EXISTS chk_projects_section,
    ALTER COLUMN section SET DEFAULT '';

