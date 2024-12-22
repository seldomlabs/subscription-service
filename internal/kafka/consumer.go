package kafka

import (
	"encoding/json"
	"log"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"subscription-service/internal/model"
)

type Consumer struct {
	consumer *kafka.Consumer
}

func NewConsumer(broker string) *Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          "subscription-service-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	return &Consumer{
		consumer: c,
	}
}

func (c *Consumer) ListenAndConsume(topic string, purchaseSubscription func(userID string, plan model.SubscriptionPlan, duration int, transactionID string) (*model.UserSubscription, error)) {
	type paymentEvent struct {
		UserID string                 `json:"user_id"`
		Plan   model.SubscriptionPlan `json:"plan"`
		Duration int				  `json:"duration"`
		TransactionID string		  `json:"transaction_id"`
	}

	err := c.consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Error subscribing to Kafka topic: %v", err)
	}

	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err == nil {
			log.Printf("Consumed message: %s", msg.Value)

			var event paymentEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}

			purchaseSubscription(event.UserID, event.Plan, event.Duration, event.TransactionID)
		} else {
			log.Printf("Error consuming message: %v", err)
		}
	}
}
