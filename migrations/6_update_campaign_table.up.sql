ALTER TABLE campaigns
DROP COLUMN amount,
ADD COLUMN streak integer NOT NULL;