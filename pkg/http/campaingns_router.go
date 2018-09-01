package http

import (
	"net/http"
	"streakr-backend/pkg/services"

	"github.com/gin-gonic/gin"
)

func CampaignRouter(data Data) {
	campaigns := data.Router.Group("/campaigns")
	{
		campaigns.POST("", func(c *gin.Context) {
			var registrationData services.NewCampaignData
			if err := c.ShouldBindJSON(&registrationData); err == nil {
				campaignId, err := registrationData.Campaign.Insert(data.Database)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, gin.H{"campaign_id": campaignId})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			}
		})

		campaigns.GET(":id", func(c *gin.Context) {
			id := c.Param("id")
			campaigns, err := services.GetCampaignById(data.Database, id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"campaigns": campaigns})
		})

		campaigns.GET("", func(c *gin.Context) {
			campaigns, err := services.GetCampaigns(data.Database)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"campaigns": campaigns})
		})

		// campaigns.POST("/:campaignId/users/:userId", func(c *gin.Context) {
		// 	campaignId := c.Param("campaignId")
		// 	userId := c.Param("userId")

		// 	streakLength, err := services.AddOrUpdateUserToCampaign(data.Database, campaignId, userId)

		// 	if err != nil {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		// 		return
		// 	}

		// 	c.JSON(http.StatusOK, gin.H{"status": streakLength})
		// })
	}
}
