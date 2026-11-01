UPDATE projects
SET plant = 'Cikarang'
WHERE plant NOT IN ('Cikarang', 'Karawang');

ALTER TABLE projects
    ALTER COLUMN plant SET DEFAULT 'Cikarang',
    ADD CONSTRAINT chk_projects_plant CHECK (plant IN ('Cikarang', 'Karawang'));

