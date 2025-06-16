package repository

import (
	"api.mijkomp.com/models/entity"
	"gorm.io/gorm"
)

type MenuRepository interface {
	Save(db *gorm.DB, payload entity.Menu) (entity.Menu, error)
	Delete(db *gorm.DB, menuId uint) error
	Search(db *gorm.DB, query *string, parentId *uint) []entity.Menu
	GetById(db *gorm.DB, menuId uint) (entity.Menu, error)

	CreateItem(db *gorm.DB, menuItem entity.MenuItem) error
	DeleteItem(db *gorm.DB, itemId uint) error
}
