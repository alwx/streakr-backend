ALTER TABLE users ADD COLUMN username TEXT NOT NULL DEFAULT "username";
ALTER TABLE users DROP COLUMN api_key;
ALTER TABLE users DROP COLUMN public_key;
ALTER TABLE users DROP COLUMN private_key;
ALTER TABLE users DROP COLUMN user_token;
ALTER TABLE users DROP COLUMN display_name;