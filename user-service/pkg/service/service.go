package service

import (
	userService "github.com/stas-bukovskiy/wish-scribe/user-service"
	"github.com/stas-bukovskiy/wish-scribe/user-service/pkg/repository"
)

type Authorization interface {
	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (userService.User, error)
}

type User interface {
	CreateUser(user userService.User) (uint, error)
	GetById(id uint) (userService.User, error)
}

type Service struct {
	Authorization
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.User),
		User:          NewUserService(repos.User),
	}
}
