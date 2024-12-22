package kafka

import (
	"encoding/json"
	"log"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Producer manages Kafka message production
type Producer struct {
	producer *kafka.Producer
}

// NewProducer initializes a Kafka producer
func NewProducer(broker string) *Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker, "debug": "all"})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	return &Producer{producer: producer}
}

// Publish sends a message to a Kafka topic
func (p *Producer) Publish(topic string, event interface{}) error {
	message, err := json.Marshal(event)
	log.Printf("Kafka event: %v", message)
	if err != nil {
		log.Printf("Failed to serialize event: %v", err)
		return fmt.Errorf("failed to produce message to Kafka: %w", err)
	}
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)

	if err != nil {
		return fmt.Errorf("failed to produce message to Kafka: %w", err)
	}
	return nil
}
