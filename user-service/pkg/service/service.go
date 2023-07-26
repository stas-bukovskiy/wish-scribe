package service

import "github.com/stas-bukovskiy/wish-scribe/user-service/pkg/repository"

type Authorization interface {
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
	return &Service{}
}
