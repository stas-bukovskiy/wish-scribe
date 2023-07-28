package repository

import (
	userService "github.com/stas-bukovskiy/wish-scribe/user-service/internal/entity"
	"gorm.io/gorm"
)

type User interface {
	CreateUser(user userService.User) (uint, error)
	GetUserByEmailAndPassword(email, password string) (userService.User, error)
	GetUserById(id uint) (userService.User, error)
}

type Repository struct {
	User
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: NewUserPostgres(db),
	}
}
