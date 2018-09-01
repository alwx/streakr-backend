package services

import (
	"database/sql"

	"github.com/google/uuid"
)

type NewCampaign struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	BadgeImageUrl string `json:"badge_image_url"`
	MinPrice      int    `json:"min_price"`
	Streak        int    `json:"streak"`
}

type Campaign struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	BadgeImageUrl string `json:"badge_image_url,omitempty"`
	MinPrice      int    `json:"min_price,omitempty"`
	Streak        int    `json:"streak,omitempty"`
}

type NewCampaignData struct {
	Campaign *NewCampaign `json:"campaign" binding:"required"`
}

func (campaign *NewCampaign) Insert(db *sql.DB) (string, error) {
	var lastInsertId string
	var err = db.QueryRow(
		"INSERT INTO campaigns(id, name, description, badge_image_url, min_price, streak) VALUES($1, $2, $3, $4, $5, $6) RETURNING id;",
		uuid.New().String(),
		campaign.Name,
		campaign.Description,
		campaign.BadgeImageUrl,
		campaign.MinPrice,
		campaign.Streak,
	).Scan(&lastInsertId)

	if err != nil {
		return "", err
	}

	return lastInsertId, nil
}

func GetCampaignById(db *sql.DB, campaignUuid string) (Campaign, error) {
	var campaign Campaign

	err := db.QueryRow(
		"SELECT id, name, description, badge_image_url, min_price, streak FROM campaigns WHERE Id = $1",
		campaignUuid,
	).Scan(&campaign.Id, &campaign.Name, &campaign.Description, &campaign.BadgeImageUrl, &campaign.MinPrice, &campaign.Streak)

	if err != nil {
		return Campaign{}, err
	}

	return campaign, nil
}

func GetCampaigns(db *sql.DB) ([]Campaign, error) {
	var campaigns []Campaign

	rows, err := db.Query("SELECT id, name, description, badge_image_url, min_price, streak FROM campaigns")
	if err != nil {
		return []Campaign{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var campaign Campaign
		err = rows.Scan(&campaign.Id, &campaign.Name, &campaign.Description, &campaign.BadgeImageUrl, &campaign.MinPrice, &campaign.Streak)
		if err == nil {
			campaigns = append(campaigns, campaign)
		}
	}

	return campaigns, nil
}
