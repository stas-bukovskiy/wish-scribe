package service

import (
	"github.com/stas-bukovskiy/wish-scribe/packages/errs"
	userService "github.com/stas-bukovskiy/wish-scribe/user-service/internal/entity"
	"github.com/stas-bukovskiy/wish-scribe/user-service/internal/repository"
	"github.com/stas-bukovskiy/wish-scribe/user-service/internal/validation"
)

type UserService struct {
	repo repository.User
}

func (s *UserService) CreateUser(user userService.User) (uint, error) {
	if !validation.IsValidEmail(user.Email) {
		return 0, errs.NewError(errs.Validation, "User email is not valid")
	}
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *UserService) GetById(id uint) (userService.User, error) {
	user, err := s.repo.GetUserById(id)
	if err != nil {
		return userService.User{}, err
	}
	return user, nil
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}
