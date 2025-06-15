package repository

import (
	"api.mijkomp.com/models/entity"
	"gorm.io/gorm"
)

type ComponentTypeRepository interface {
	Save(db *gorm.DB, payload entity.ComponentType) (entity.ComponentType, error)
	Delete(db *gorm.DB, componentTypeId uint) error
	Search(db *gorm.DB, query *string) []entity.ComponentType
	GetById(db *gorm.DB, componentTypeId uint) (entity.ComponentType, error)
}
