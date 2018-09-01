ALTER TABLE campaigns
DROP COLUMN streak,
ADD COLUMN amount integer NOT NULL;