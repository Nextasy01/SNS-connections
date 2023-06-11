package entity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID `gorm:"type:varchar(16);primaryKey"`
	Email      string    `json:"email" binding:"required,email" gorm:"type:varchar(256);not null;unique"`
	Username   string    `json:"username" binding:"required" gorm:"type:varchar(32);not null;unique"`
	ProfilePic string    `gorm:"type:string"`
	Password   string    `json:"password" binding:"required" gorm:"type:varchar(32);not null"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func NewUser() *User {
	return &User{}
}
func (u *User) BeforeCreate(db *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	return nil
}

func (u *User) PrepareGive() {
	u.Password = ""
}
