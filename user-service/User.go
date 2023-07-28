package user_service

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" binging:"required"`
	Email     string `json:"email" binging:"required" gorm:"uniqueIndex"`
	Password  string `json:"-" gorm:"column:password_hash"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
