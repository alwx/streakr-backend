ALTER TABLE campaigns
ADD COLUMN badge_image_url TEXT NOT NULL,
DROP COLUMN prize_description;