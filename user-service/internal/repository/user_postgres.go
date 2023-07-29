package repository

import (
	"errors"
	"github.com/stas-bukovskiy/wish-scribe/packages/errs"
	"github.com/stas-bukovskiy/wish-scribe/packages/logger"
	userService "github.com/stas-bukovskiy/wish-scribe/user-service/internal/entity"
	"gorm.io/gorm"
)

type UserPostgres struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewUserPostgres(db *gorm.DB, logger logger.Logger) *UserPostgres {
	return &UserPostgres{db: db, logger: logger}
}

func (us *UserPostgres) CreateUser(user userService.User) (uint, error) {
	log := us.logger.Named("Create").With("user", user)
	tx := us.db.Begin()
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
		return 0, errs.NewError(errs.AlreadyExist, "user with such mail already exists")
	}

	result := tx.Create(&user)

	tx.Commit()
	log.Debug("successfully created user")
	return user.ID, result.Error
}

func (us *UserPostgres) GetUserByEmailAndPassword(email, password string) (userService.User, error) {
	log := us.logger.Named("GetUserByEmailAndPassword").With("email", email, "password", password)
	var user userService.User
	err := us.db.Where("email = ? AND password_hash = ?", email, password).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errs.NewError(errs.NotFound, "User not found")
	}
	if err != nil {
		return user, err
	}
	log.Debug("successfully found user")
	return user, nil
}

func (us *UserPostgres) GetUserById(id uint) (userService.User, error) {
	log := us.logger.Named("GetUserByEmailAndPassword").With("id", id)
	var user userService.User
	err := us.db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errs.NewError(errs.NotFound, "User not found")
	}
	if err != nil {
		return user, err
	}
	log.Debug("successfully found user")
	return user, nil
}
