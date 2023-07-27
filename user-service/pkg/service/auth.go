package service

import (
	"crypto/sha1"
	"fmt"
	userService "github.com/stas-bukovskiy/wish-scribe/user-service"
	"github.com/stas-bukovskiy/wish-scribe/user-service/pkg/errs"
	"github.com/stas-bukovskiy/wish-scribe/user-service/pkg/repository"
	"github.com/stas-bukovskiy/wish-scribe/user-service/pkg/utils"
)

const salt = "jevBH89BC9cbdsc298dUCXbzasxOZox"

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user userService.User) (uint, error) {
	if !utils.IsValidEmail(user.Email) {
		return 0, errs.NewError(errs.Validation, "User email is not valid")
	}
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
