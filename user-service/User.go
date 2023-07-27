package user_service

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" binging:"required"`
	Email     string    `json:"email" binging:"required" gorm:"uniqueIndex"`
	Password  string    `json:"-" gorm:"column:password_hash"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
