package repository

import (
	"api.mijkomp.com/models/entity"
	"gorm.io/gorm"
)

type ProductCategoryRepository interface {
	Save(db *gorm.DB, payload entity.ProductCategory) (entity.ProductCategory, error)
	Delete(db *gorm.DB, categoryId uint) error
	Search(db *gorm.DB, query *string) []entity.ProductCategory
	GetById(db *gorm.DB, categoryId uint) (entity.ProductCategory, error)
}
