ALTER TABLE campaign_user
ADD CONSTRAINT FK_users_id FOREIGN KEY (userId)
    REFERENCES users(id)
    ON DELETE CASCADE;

ALTER TABLE campaign_user
ADD CONSTRAINT FK_campaign_id FOREIGN KEY (campaignId)
    REFERENCES campaigns(id)
    ON DELETE CASCADE;