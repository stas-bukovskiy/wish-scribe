package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/stas-bukovskiy/wish-scribe/packages/errs"
	userService "github.com/stas-bukovskiy/wish-scribe/user-service/internal/entity"
	"github.com/stas-bukovskiy/wish-scribe/user-service/internal/repository"
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
	repo repository.User
}

func NewAuthService(repo repository.User) *AuthService {
	return &AuthService{repo: repo}
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) GenerateToken(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmailAndPassword(email, generatePasswordHash(password))
	if err != nil {
		if errs.KindIs(errs.NotFound, err) {
			return "", errs.NewError(errs.Unauthorized, "Invalid email or password")
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

func (s *AuthService) ParseToken(tokenToParse string) (userService.User, error) {
	token, err := jwt.ParseWithClaims(tokenToParse, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.NewError(errs.Unauthorized, "Invalid signing method")
		}

		return []byte(singingKey), nil
	})
	if err != nil {
		return userService.User{}, errs.NewError(errs.Unauthorized, "Invalid access token")
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return userService.User{}, errs.NewError(errs.Unauthorized, "Invalid access token")
	}
	userId := claims.UserId
	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return userService.User{}, errs.NewError(errs.Unauthorized, "Invalid access token")
	}
	return user, nil
}
