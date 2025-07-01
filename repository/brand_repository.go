package repository

import (
	"api.mijkomp.com/models/entity"
	"gorm.io/gorm"
)

type BrandRepository interface {
	Save(db *gorm.DB, payload entity.Brand) (entity.Brand, error)
	Delete(db *gorm.DB, brandId uint) error
	Search(db *gorm.DB, query *string) []entity.Brand
	GetById(db *gorm.DB, brandId uint) (entity.Brand, error)
}
