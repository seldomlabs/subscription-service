package repository

import (
	"fmt"
	"subscription-service/internal/model"
	"errors"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("record not found")

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Get(userID string) (*model.UserSubscription, error) {
	var subscription model.UserSubscription
	if err := r.db.Where("user_id = ?", userID).First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("subscription not found for user_id: %s", userID)
		}
		return nil, err
	}
	return &subscription, nil
}

func (r *SubscriptionRepository) Save(subscription *model.UserSubscription) error {
	if err := r.db.Save(subscription).Error; err != nil {
		return fmt.Errorf("failed to save subscription: %w", err)
	}
	return nil
}
