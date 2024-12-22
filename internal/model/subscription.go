package model

import "time"

type SubscriptionPlan string

const (
	Standard SubscriptionPlan = "Standard"
	Gold     SubscriptionPlan = "Gold"
	Platinum SubscriptionPlan = "Platinum"
)

type SubscriptionStatus string

const (
	Active     SubscriptionStatus = "Active"
	Canceled   SubscriptionStatus = "Canceled"
	OnHold     SubscriptionStatus = "OnHold"
	Expired    SubscriptionStatus = "Expired"
)

type UserSubscription struct {
	ID               uint               `gorm:"primaryKey" json:"id"`
	UserID           string             `gorm:"type:varchar(255);uniqueIndex;not null" json:"user_id"` 
	Plan             SubscriptionPlan   `gorm:"type:varchar(255);not null" json:"plan"`
	Status           SubscriptionStatus `gorm:"type:varchar(255)" json:"status"`
	AutoRenew        bool               `gorm:"default:true" json:"auto_renew"`  
	SubscriptionStartDate  time.Time    `gorm:"not null" json:"subscription_start_date"` 
	SubscriptionEndDate  time.Time      `gorm:"not null" json:"subscription_end_date"` 
	TransactionID      string           `gorm:"type:varchar(255);not null" json:"transaction_id"` 
	CreatedAt        time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
}
