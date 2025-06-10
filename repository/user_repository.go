package repository

import (
	"api.mijkomp.com/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(db *gorm.DB, user entity.User) (entity.User, error)
	GetById(db *gorm.DB, userId uint) (entity.User, error)
	GetByEmail(db *gorm.DB, email string) (entity.User, error)
	SaveToken(db *gorm.DB, token entity.UserToken) error
	HasToken(db *gorm.DB, userId uint, tokenId uuid.UUID) bool
	CreateUserVerification(db *gorm.DB, userVerification entity.UserVerification) error
	GetUserVerification(db *gorm.DB, userId uint, code string) (entity.UserVerification, error)
}
