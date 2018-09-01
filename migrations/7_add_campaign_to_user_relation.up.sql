CREATE TABLE IF NOT EXISTS campaign_user (
 campaignId UUID NOT NULL,
 userId UUID NOT NULL,
 streak_length integer NOT NULL DEFAULT 0
)