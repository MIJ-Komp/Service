package repository_impl

import (
	"api.mijkomp.com/models/entity"
	"gorm.io/gorm"
)

type ComponentTypeRepositoryImpl struct {
}

func NewComponentTypeRepository() *ComponentTypeRepositoryImpl {
	return &ComponentTypeRepositoryImpl{}
}

func (repository *ComponentTypeRepositoryImpl) Save(db *gorm.DB, componentType entity.ComponentType) (entity.ComponentType, error) {
	err := db.Save(&componentType).Error
	return componentType, err
}

func (repository *ComponentTypeRepositoryImpl) Delete(db *gorm.DB, componentTypeId uint) error {
	err := db.Where("id = ?", componentTypeId).Delete(&entity.ComponentType{}).Error
	return err
}

func (repository *ComponentTypeRepositoryImpl) Search(db *gorm.DB, query *string) []entity.ComponentType {

	productCategories := []entity.ComponentType{}

	queries := db.Model(&entity.ComponentType{})
	if query != nil {
		q := "%" + *query + "%"
		queries.Where("CONCAT(code, ' ', name, ' ', description) ILIKE ?", q)
	}

	queries.Order("component_types.modified_at desc").Find(&productCategories)

	return productCategories
}

func (repository *ComponentTypeRepositoryImpl) GetById(db *gorm.DB, componentTypeId uint) (entity.ComponentType, error) {
	id := componentTypeId
	var componentTypes entity.ComponentType
	err := db.
		Preload("CreatedBy").
		Preload("ModifiedBy").
		First(&componentTypes, "id = ?", id).Error

	return componentTypes, err
}
