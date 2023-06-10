package entity

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID           uuid.UUID `gorm:"type:varchar(16);primaryKey"`
	Title        string    `gorm:"type:string; not null"`
	DriveId      string    `gorm:"type:string"`
	VideoId      string    `gorm:"type:string; not null"`
	PlatformFrom string    `gorm:"type:string; not null"`
	PlatformTo   string    `gorm:"type:string; not null"`
	Creator      User      `gorm:"constraint:OnUpdate:RESTRICT,onDelete:RESTRICT;"`
	CreatorId    uuid.UUID
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func NewPost() *Post {
	return &Post{}
}
