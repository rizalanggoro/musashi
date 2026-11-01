ALTER TABLE projects
    DROP COLUMN IF EXISTS progress_xxx,
    ADD COLUMN IF NOT EXISTS progress_ppt text NOT NULL DEFAULT 'pending' CHECK (progress_ppt IN ('pending', 'done')),
    ADD COLUMN IF NOT EXISTS progress_drawing text NOT NULL DEFAULT 'pending' CHECK (progress_drawing IN ('pending', 'done')),
    ADD COLUMN IF NOT EXISTS progress_vendor text NOT NULL DEFAULT 'pending' CHECK (progress_vendor IN ('pending', 'done')),
    ADD COLUMN IF NOT EXISTS progress_pr text NOT NULL DEFAULT 'pending' CHECK (progress_pr IN ('pending', 'done')),
    ADD COLUMN IF NOT EXISTS progress_ringi text NOT NULL DEFAULT 'pending' CHECK (progress_ringi IN ('pending', 'done')),
    ADD COLUMN IF NOT EXISTS progress_po text NOT NULL DEFAULT 'pending' CHECK (progress_po IN ('pending', 'done')),
    ADD COLUMN IF NOT EXISTS progress_construction text NOT NULL DEFAULT 'pending' CHECK (progress_construction IN ('pending', 'done')),
    ADD COLUMN IF NOT EXISTS progress_bast text NOT NULL DEFAULT 'pending' CHECK (progress_bast IN ('pending', 'done'));

