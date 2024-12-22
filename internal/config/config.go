package config

import (
	"os"
)

type Config struct {
	ServerPort  string
	KafkaBroker string
	SubscriptionTopic string
	PaymentTopic string
	MySQLDSN string
}

func Load() *Config {
	serverPort := getEnv("SERVER_PORT", "5002")
	kafkaBroker := getEnv("KAFKA_BROKER", "34.227.32.140:9092")
	subscriptionTopic := getEnv("KAFKA_SUBSCRIPTION_TOPIC", "subscription")
	paymentTopic := getEnv("KAFKA_PAYMENT_TOPIC", "payment")
	MySQLDSN := getEnv("MySQLDSN", "root:root_password@tcp(98.83.138.170:3306)/projectx?parseTime=true")

	return &Config{
		ServerPort:  serverPort,
		KafkaBroker: kafkaBroker,
		SubscriptionTopic: subscriptionTopic,
		PaymentTopic: paymentTopic,
		MySQLDSN: MySQLDSN,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
