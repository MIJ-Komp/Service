package repository_impl

import (
	"api.mijkomp.com/models/entity"
	"gorm.io/gorm"
)

type MenuRepositoryImpl struct {
}

func NewMenuRepository() *MenuRepositoryImpl {
	return &MenuRepositoryImpl{}
}

func (repository *MenuRepositoryImpl) Save(db *gorm.DB, menu entity.Menu) (entity.Menu, error) {
	err := db.Save(&menu).Error
	return menu, err
}

func (repository *MenuRepositoryImpl) Delete(db *gorm.DB, menuId uint) error {
	err := db.Where("id = ?", menuId).Delete(&entity.Menu{}).Error
	return err
}

func (repository *MenuRepositoryImpl) Search(db *gorm.DB, query *string, parentId *uint) []entity.Menu {

	menus := []entity.Menu{}

	queries := db.
		Preload("MenuItems").
		Preload("MenuItems.ProductCategory").
		Model(&entity.Menu{})

	if query != nil {
		queries.Where("name like ?", "%"+*query+"%")
	}

	if parentId != nil {
		queries.Where("parent_id = ?", parentId)
	}
	queries.Order("name desc").Find(&menus)

	return menus
}

func (repository *MenuRepositoryImpl) GetById(db *gorm.DB, menuId uint) (entity.Menu, error) {
	id := menuId
	var menus entity.Menu
	err := db.
		Preload("CreatedBy").
		Preload("ModifiedBy").
		First(&menus, "id = ?", id).Error

	return menus, err
}

func (repository *MenuRepositoryImpl) CreateItem(db *gorm.DB, menuItem entity.MenuItem) error {
	err := db.Save(&menuItem).Error

	return err
}

func (repository *MenuRepositoryImpl) DeleteItem(db *gorm.DB, itemId uint) error {
	err := db.Where("id = ?", itemId).Delete(&entity.MenuItem{}).Error
	return err
}
