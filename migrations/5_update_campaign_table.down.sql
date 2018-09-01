ALTER TABLE campaigns
DROP COLUMN badge_image_url,
ADD COLUMN prize_description text NOT NULL;