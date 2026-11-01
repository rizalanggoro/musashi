ALTER TABLE projects
    DROP COLUMN IF EXISTS progress_ppt,
    DROP COLUMN IF EXISTS progress_drawing,
    DROP COLUMN IF EXISTS progress_vendor,
    DROP COLUMN IF EXISTS progress_pr,
    DROP COLUMN IF EXISTS progress_ringi,
    DROP COLUMN IF EXISTS progress_po,
    DROP COLUMN IF EXISTS progress_construction,
    DROP COLUMN IF EXISTS progress_bast,
    ADD COLUMN IF NOT EXISTS progress_xxx text NOT NULL DEFAULT 'pending' CHECK (progress_xxx IN ('pending', 'done'));

