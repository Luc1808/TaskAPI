CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- tasks table
CREATE TABLE IF NOT EXISTS public.tasks (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	title TEXT NOT NULL
		CHECK (char_length(title) BETWEEN 1 AND 200),
	description TEXT,
	status TEXT NOT NULL DEFAULT 'todo' -- Could also have used VARCHAR; Just to be different
		CHECK (status IN ('todo', 'progress', 'done')),
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- auto-update "update_at" on row UPDATE
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
	NEW.update_at = now();
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_tasks_set_updated_at ON public.tasks;
CREATE TRIGGER trg_tasks_set_updated_at
BEFORE UPDATE ON public.tasks
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Helpful composite index for list filters
CREATE INDEX IF NOT EXISTS idx_tasks_status_created_at
ON public.tasks (status, created_at DESC);
