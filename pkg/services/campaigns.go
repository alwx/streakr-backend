package services

import (
	"database/sql"

	"github.com/google/uuid"
)

type NewCampaign struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	PrizeDescription string `json:"prize_description"`
	MinPrice         int    `json:"min_price"`
	Amount           int    `json:"amount"`
}

type Campaign struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	PrizeDescription string `json:"prize_description,omitempty"`
	MinPrice         int    `json:"min_price,omitempty"`
	Amount           int    `json:"amount,omitempty"`
}

type NewCampaignData struct {
	Campaign *NewCampaign `json:"campaign" binding:"required"`
}

func (campaign *NewCampaign) Insert(db *sql.DB) (string, error) {
	var lastInsertId string
	var err = db.QueryRow(
		"INSERT INTO campaigns(id, name, description, prize_description, min_price, amount) VALUES($1, $2, $3, $4, $5, $6) RETURNING id;",
		uuid.New().String(),
		campaign.Name,
		campaign.Description,
		campaign.PrizeDescription,
		campaign.MinPrice,
		campaign.Amount,
	).Scan(&lastInsertId)

	if err != nil {
		return "", err
	}

	return lastInsertId, nil
}

func GetCampaignById(db *sql.DB, campaignUuid string) (Campaign, error) {
	var campaign Campaign

	err := db.QueryRow(
		"SELECT id, name, description, prize_description, min_price, amount FROM campaigns WHERE Id = $1",
		campaignUuid,
	).Scan(&campaign.Id, &campaign.Name, &campaign.Description, &campaign.PrizeDescription, &campaign.MinPrice, &campaign.Amount)

	if err != nil {
		return Campaign{}, err
	}

	return campaign, nil
}

func GetCampaigns(db *sql.DB) ([]Campaign, error) {
	var campaigns []Campaign

	rows, err := db.Query("SELECT id, name, description FROM campaigns")
	if err != nil {
		return []Campaign{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var campaign Campaign
		err = rows.Scan(&campaign.Id, &campaign.Name, &campaign.Description)
		if err == nil {
			campaigns = append(campaigns, campaign)
		}
	}

	return campaigns, nil
}
