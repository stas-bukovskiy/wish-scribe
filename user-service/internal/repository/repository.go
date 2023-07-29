package repository

import (
	"github.com/stas-bukovskiy/wish-scribe/packages/logger"
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

func NewRepository(db *gorm.DB, logger logger.Logger) *Repository {
	return &Repository{
		User: NewUserPostgres(db, logger),
	}
}
