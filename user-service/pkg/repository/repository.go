package repository

import (
	userService "github.com/stas-bukovskiy/wish-scribe/user-service"
	"gorm.io/gorm"
)

type Authorization interface {
	CreateUser(user userService.User) (uint, error)
	GetUserByEmailAndPassword(email, password string) (userService.User, error)
}

type Token interface {
}

type User interface {
}

type Repository struct {
	Authorization
	Token
	User
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
