package repository_impl

import (
	"errors"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Save(db *gorm.DB, user entity.User) (entity.User, error) {
	err := db.Save(&user).Error
	return user, err
}

func (repository *UserRepositoryImpl) GetById(db *gorm.DB, userId uint) (entity.User, error) {
	id := userId
	var users entity.User
	err := db.First(&users, "id = ?", id).Error

	return users, err
}

func (repository *UserRepositoryImpl) GetByEmail(db *gorm.DB, email string) (entity.User, error) {
	var user entity.User
	err := db.First(&user, "email = ?", email).Error
	return user, err
}

func (repository *UserRepositoryImpl) SaveToken(db *gorm.DB, token entity.UserToken) error {

	if err := db.Save(&token).Error; err != nil {
		exception.PanicIfNeeded(err)
	}
	return nil
}

func (repository *UserRepositoryImpl) HasToken(db *gorm.DB, userId uint, tokenId uuid.UUID) bool {

	var userToken entity.UserToken

	if err := db.First(&userToken, "user_id = ? and token = ?", userId, tokenId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false // User does not exist
		}
		return false
	}
	return true
}

func (repository *UserRepositoryImpl) CreateUserVerification(db *gorm.DB, userVerification entity.UserVerification) error {
	err := db.Save(&userVerification).Error
	return err
}

func (repository *UserRepositoryImpl) GetUserVerification(db *gorm.DB, userId uint, code string) (entity.UserVerification, error) {
	var userVerification entity.UserVerification

	err := db.First(&userVerification, "user_id = ? and code = ?", userId, code).Error

	return userVerification, err
}
