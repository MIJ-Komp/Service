package repository_impl

import (
	"errors"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers/logger"
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
	logger.LogDBOperation("Save", "Save user", user.Email)
	err := db.Save(&user).Error
	if err != nil {
		logger.LogError(err.Error())
	}
	return user, err
}

func (repository *UserRepositoryImpl) GetById(db *gorm.DB, userId uint) (entity.User, error) {
	logger.LogDBOperation("GetById", "Get user by ID", userId)
	id := userId
	var users entity.User
	err := db.First(&users, "id = ?", id).Error
	if err != nil {
		logger.LogError(err.Error())
	}
	return users, err
}

func (repository *UserRepositoryImpl) GetByEmail(db *gorm.DB, email string) (entity.User, error) {
	logger.LogDBOperation("GetByEmail", "Get user by email", email)
	var user entity.User
	err := db.First(&user, "email = ?", email).Error
	if err != nil {
		logger.LogError(err.Error())
	}
	return user, err
}

func (repository *UserRepositoryImpl) SaveToken(db *gorm.DB, token entity.UserToken) error {
	logger.LogDBOperation("SaveToken", "Save user token", token.UserId)
	if err := db.Save(&token).Error; err != nil {
		logger.LogError(err.Error())
		exception.PanicIfNeeded(err)
	}
	return nil
}

func (repository *UserRepositoryImpl) HasToken(db *gorm.DB, userId uint, tokenId uuid.UUID) bool {
	logger.LogDBOperation("HasToken", "Check user token", userId, tokenId)
	var userToken entity.UserToken

	if err := db.First(&userToken, "user_id = ? and token = ?", userId, tokenId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogInfo("Token not found for user")
			return false // User does not exist
		}
		logger.LogError(err.Error())
		return false
	}
	return true
}

func (repository *UserRepositoryImpl) CreateUserVerification(db *gorm.DB, userVerification entity.UserVerification) error {
	logger.LogDBOperation("CreateUserVerification", "Create user verification", userVerification.UserId)
	err := db.Save(&userVerification).Error
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (repository *UserRepositoryImpl) GetUserVerification(db *gorm.DB, userId uint, code string) (entity.UserVerification, error) {
	logger.LogDBOperation("GetUserVerification", "Get user verification", userId, code)
	var userVerification entity.UserVerification

	err := db.First(&userVerification, "user_id = ? and code = ?", userId, code).Error
	if err != nil {
		logger.LogError(err.Error())
	}

	return userVerification, err
}
