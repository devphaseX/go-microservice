-- Drop the trigger
DROP TRIGGER IF EXISTS update_users_updated_at ON public.users;

-- Drop the function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop the index
DROP INDEX IF EXISTS idx_users_email;

-- Drop the users table
DROP TABLE IF EXISTS public.users;
