ALTER TABLE users DROP COLUMN username;
ALTER TABLE users ADD COLUMN api_key TEXT;
ALTER TABLE users ADD COLUMN public_key TEXT;
ALTER TABLE users ADD COLUMN private_key TEXT;
ALTER TABLE users ADD COLUMN user_token TEXT;
ALTER TABLE users ADD COLUMN display_name TEXT;