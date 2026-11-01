UPDATE projects
SET section = 'GA Building'
WHERE section NOT IN ('GA Building', 'EHS', 'Wellfare');

ALTER TABLE projects
    ALTER COLUMN section SET DEFAULT 'GA Building',
    ADD CONSTRAINT chk_projects_section CHECK (section IN ('GA Building', 'EHS', 'Wellfare'));

