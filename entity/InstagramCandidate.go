package entity

import (
	"time"

	"github.com/google/uuid"
)

type InstagramCandidate struct {
	ID        uuid.UUID        `gorm:"type:varchar(16);primaryKey"`
	Caption   string           `gorm:"type:string"`
	MediaUrl  string           `gorm:"type:string"`
	Permalink string           `gorm:"type:string"`
	ShortCode string           `gorm:"type:string"`
	VideoId   string           `gorm:"type:string; not null"`
	Creator   InstagramAccount `gorm:"constraint:OnUpdate:RESTRICT,onDelete:CASCADE;"`
	CreatorId uuid.UUID
	Timestamp time.Time `gorm:"type:datetime; not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func NewInstagramCandidate() *InstagramCandidate {
	return &InstagramCandidate{}
}
