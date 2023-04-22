package entity

import (
	"time"

	"github.com/google/uuid"
)

type InstagramAccount struct {
	ID                  uuid.UUID `gorm:"type:varchar(16);primaryKey"`
	Username            string    `gorm:"type:varchar(32);not null"`
	InstagramPrivateID  string    `gorm:"type:string"`
	InstagramBusinessID string    `gorm:"type:string"`
	InstagramUserID     string    `gorm:"type:string"`
	ProfilePic          string    `gorm:"type:string"`
	AccessToken         string    `gorm:"type:string"`
	TokenType           string    `gorm:"type:string"`
	User                User      `gorm:"constraint:OnUpdate:RESTRICT,onDelete:CASCADE;"`
	UserID              uuid.UUID
	ExpiresAt           time.Time `gorm:"column:expiration_date"`
	CreatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func NewInstagramAccount() *InstagramAccount {
	return &InstagramAccount{}
}
