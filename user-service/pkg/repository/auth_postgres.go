package repository

import (
	"errors"
	userService "github.com/stas-bukovskiy/wish-scribe/user-service"
	"github.com/stas-bukovskiy/wish-scribe/user-service/pkg/errs"
	"gorm.io/gorm"
)

type AuthPostgres struct {
	db *gorm.DB
}

func NewAuthPostgres(db *gorm.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (ap *AuthPostgres) CreateUser(user userService.User) (uint, error) {
	tx := ap.db.Begin()
	var exists bool
	err := tx.Model(&userService.User{}).Select("count(*) > 0").
		Where("email = ?", user.Email).
		Find(&exists).
		Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if exists {
		tx.Rollback()
		return 0, errs.NewError(errs.Exist, "user with such mail already exists")
	}

	result := tx.Create(&user)

	tx.Commit()
	return user.ID, result.Error
}

func (ap *AuthPostgres) GetUserByEmailAndPassword(email, password string) (userService.User, error) {
	var user userService.User
	err := ap.db.Where("email = ? AND password_hash = ?", email, password).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errs.NewError(errs.NotExist, "User not found")
	}
	if err != nil {
		return user, err
	}
	return user, nil
}

func (ap *AuthPostgres) GetUserById(id uint) (userService.User, error) {
	var user userService.User
	err := ap.db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errs.NewError(errs.NotExist, "User not found")
	}
	if err != nil {
		return user, err
	}
	return user, nil
}
