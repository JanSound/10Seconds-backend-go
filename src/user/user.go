package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	Name      string
	Email     string
	createdAt time.Time
	UpdatedAt time.Time
}
