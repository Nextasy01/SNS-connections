package entity

import (
	"time"

	"github.com/google/uuid"
)

type YoutubeCandidate struct {
	ID          uuid.UUID     `gorm:"type:varchar(16);primaryKey"`
	Title       string        `gorm:"type:string; not null"`
	Description string        `gorm:"type:string"`
	ChannelId   string        `gorm:"type:string; not null"`
	PublishedAt time.Time     `gorm:"type:datetime; not null"`
	VideoId     string        `gorm:"type:string; not null"`
	Creator     GoogleAccount `gorm:"constraint:OnUpdate:RESTRICT,onDelete:CASCADE;"`
	CreatorId   uuid.UUID
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func NewYouTubeCandidate() *YoutubeCandidate {
	return &YoutubeCandidate{}
}
