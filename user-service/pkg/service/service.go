package service

import (
	userService "github.com/stas-bukovskiy/wish-scribe/user-service"
	"github.com/stas-bukovskiy/wish-scribe/user-service/pkg/repository"
)

type Authorization interface {
	CreateUser(user userService.User) (uint, error)
}

type Token interface {
}

type User interface {
}

type Service struct {
	Authorization
	Token
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
