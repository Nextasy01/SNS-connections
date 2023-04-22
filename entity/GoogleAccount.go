package entity

import (
	"time"

	"github.com/google/uuid"
)

type GoogleAccount struct {
	ID           uuid.UUID `gorm:"type:varchar(16);primaryKey"`
	Email        string    `gorm:"type:varchar(256);not null;unique"`
	Username     string    `gorm:"type:varchar(32);not null"`
	ProfilePic   string    `gorm:"type:string"`
	RefreshToken string    `gorm:"type:string"`
	AccessToken  string    `gorm:"type:string"`
	TokenType    string    `gorm:"type:string"`
	User         User      `gorm:"constraint:OnUpdate:RESTRICT,onDelete:CASCADE;"`
	UserID       uuid.UUID
	ExpiresAt    time.Time `gorm:"column:expiration_date"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func NewGoogleAccount() *GoogleAccount {
	return &GoogleAccount{}
}
