package service

import (
	"errors"
	"time"
	"log"
	"subscription-service/internal/kafka"
	"subscription-service/internal/model"
	"subscription-service/internal/repository"
)

type SubscriptionService struct {
	repo     *repository.SubscriptionRepository
	producer *kafka.Producer
	subscriptionTopic string
}

func NewSubscriptionService(repo *repository.SubscriptionRepository, producer *kafka.Producer, topic string) *SubscriptionService {
	log.Printf("Topic=%s", topic)
	return &SubscriptionService{
		repo:     repo,
		producer: producer,
		subscriptionTopic: topic,
	}
}

func (s *SubscriptionService) PurchaseSubscription(userID string, plan model.SubscriptionPlan, duration int, transactionID string) (*model.UserSubscription, error) {
	sub := &model.UserSubscription{
		UserID:           userID,
		Plan:             plan,
		Status:           model.Active,
		AutoRenew:        true,
		SubscriptionStartDate:  time.Now(), 
		SubscriptionEndDate:  time.Now().AddDate(0, 0, duration),
		TransactionID: transactionID,
	}

	if err := s.repo.Save(sub); err != nil {
		return nil, err
	}

	log.Printf("Kafka Topic=%s", s.subscriptionTopic)
	if err := s.producer.Publish(s.subscriptionTopic, map[string]interface{}{
		"event": "PURCHASED",
		"data":  sub,
	}); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *SubscriptionService) UpdateSubscription(sub *model.UserSubscription) (*model.UserSubscription, error) {

	if err := s.repo.Save(sub); err != nil {
		return nil, err
	}

	if err := s.producer.Publish(s.subscriptionTopic, map[string]interface{}{
		"event": "SWITCH",
		"data":  sub,
	}); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *SubscriptionService) GetSubscription(userID string) (*model.UserSubscription, error) {
	sub, err := s.repo.Get(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errors.New("subscription not found")
		}
		return nil, err
	}

	return sub, nil
}

func (s *SubscriptionService) CancelSubscription(userID string) error {
	sub, err := s.repo.Get(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return errors.New("subscription not found")
		}
		return err
	}

	// Update subscription status to Canceled
	sub.Status = model.Canceled
	sub.AutoRenew = false

	// Save changes to the repository
	if err := s.repo.Save(sub); err != nil {
		return err
	}

	// Publish Kafka event
	return s.producer.Publish(s.subscriptionTopic, map[string]interface{}{
		"event": "CANCELED",
		"data":  sub,
	})
}
