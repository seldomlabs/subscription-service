package main

import (
	"log"
	"subscription-service/internal/api"
	"subscription-service/internal/config"
	"subscription-service/internal/kafka"
	"subscription-service/internal/repository"
	"subscription-service/internal/service"
	"subscription-service/internal/model"
	"gorm.io/driver/mysql"
    "gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to MySQL: %v", err)
    }
    db.AutoMigrate(&model.UserSubscription{})

	kafkaProducer := kafka.NewProducer(cfg.KafkaBroker)

	subRepo := repository.NewSubscriptionRepository(db)

	subService := service.NewSubscriptionService(subRepo, kafkaProducer, cfg.SubscriptionTopic)

	kafkaConsumer := kafka.NewConsumer(cfg.KafkaBroker)

	go kafkaConsumer.ListenAndConsume(cfg.PaymentTopic,subService.PurchaseSubscription)

	// Set up HTTP API
	router := gin.Default()
	api.RegisterRoutes(router, subService)

	// Start the server
	log.Printf("Server starting on port %s...", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
