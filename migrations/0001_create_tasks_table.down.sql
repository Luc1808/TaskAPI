-- drop in reverse order of dependencies
DROP TRIGGER IF EXISTS trg_tasks_set_updated_at ON public.tasks;
DROP FUNCTION IF EXISTS set_updated_at();

DROP INDEX IF EXISTS idx_tasks_status_created_at;

DROP TABLE IF EXISTS public.tasks;


