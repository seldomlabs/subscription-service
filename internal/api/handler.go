package api

import (
	"net/http"
	"subscription-service/internal/model"
	"subscription-service/internal/service"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up API routes
func RegisterRoutes(router *gin.Engine, svc *service.SubscriptionService) {
	api := router.Group("/subscriptions")
	{
		api.POST("/purchase", func(c *gin.Context) {
			var req struct {
				UserID string                 `json:"user_id"`
				Plan   model.SubscriptionPlan `json:"plan"`
				Duration int				  `json:"duration"`
				TransactionID string		  `json:"transaction_id"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			sub,err := svc.PurchaseSubscription(req.UserID, req.Plan, req.Duration, req.TransactionID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, sub)
		})

		api.PATCH("/update", func(c *gin.Context) {
			var req struct {
				UserID        string                  `json:"user_id" binding:"required"` 
				Plan          *model.SubscriptionPlan `json:"plan"`                      
				Duration      *int                    `json:"duration"`                
				TransactionID *string                 `json:"transaction_id"`           
			}
		
			// Bind and validate JSON input
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required and must be valid"})
				return
			}
		
			existingSub, err := svc.GetSubscription(req.UserID) 
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
				return
			}

			if req.Plan != nil {
				existingSub.Plan = *req.Plan
			}
			if req.Duration != nil {
				existingSub.SubscriptionEndDate = existingSub.SubscriptionStartDate.AddDate(0,0,*req.Duration)
			}
			if req.TransactionID != nil {
				existingSub.TransactionID = *req.TransactionID
			}
		
			updatedSub, err := svc.UpdateSubscription(existingSub)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		
			c.JSON(http.StatusOK, updatedSub)
		})


		api.GET("/:user_id", func(c *gin.Context) {
			userID := c.Param("user_id")
			sub, err := svc.GetSubscription(userID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, sub)
		})

		api.DELETE("/:user_id", func(c *gin.Context) {
			userID := c.Param("user_id")
			err := svc.CancelSubscription(userID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Subscription canceled"})
		})
	}
}
