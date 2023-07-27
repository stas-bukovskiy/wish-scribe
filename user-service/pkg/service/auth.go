package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	userService "github.com/stas-bukovskiy/wish-scribe/user-service"
	"github.com/stas-bukovskiy/wish-scribe/user-service/pkg/errs"
	"github.com/stas-bukovskiy/wish-scribe/user-service/pkg/repository"
	"github.com/stas-bukovskiy/wish-scribe/user-service/pkg/utils"
	"time"
)

const (
	salt       = "jevBH89BC9cbdsc298dUCXbzasxOZox"
	singingKey = "pISc0ODSDC9023onc90132sdcu19DUISuisdUSDcx"
	tokenTTL   = 24 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId uint `json:"user_id"`
}
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
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) GenerateToken(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmailAndPassword(email, generatePasswordHash(password))
	if err != nil {
		if errs.KindIs(errs.NotExist, err) {
			return "", errs.NewError(errs.Unanticipated, "Invalid email or password")
		}
		return "", err
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(tokenTTL).Unix(),
			IssuedAt:  now.Unix(),
			Subject:   user.Email,
		},
		UserId: user.ID,
	})

	return token.SignedString([]byte(singingKey))
}
