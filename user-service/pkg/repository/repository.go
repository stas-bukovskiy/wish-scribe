package repository

import "gorm.io/gorm"

type Authorization interface {
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
	return &Repository{}
}
