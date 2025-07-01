package repository_impl

import (
	"api.mijkomp.com/models/entity"
	"gorm.io/gorm"
)

type BrandRepositoryImpl struct {
}

func NewBrandRepository() *BrandRepositoryImpl {
	return &BrandRepositoryImpl{}
}

func (repository *BrandRepositoryImpl) Save(db *gorm.DB, brand entity.Brand) (entity.Brand, error) {
	err := db.Save(&brand).Error
	return brand, err
}

func (repository *BrandRepositoryImpl) Delete(db *gorm.DB, brandId uint) error {
	err := db.Where("id = ?", brandId).Delete(&entity.Brand{}).Error
	return err
}

func (repository *BrandRepositoryImpl) Search(db *gorm.DB, query *string) []entity.Brand {

	productCategories := []entity.Brand{}

	queries := db.Model(&entity.Brand{})

	if query != nil {
		queries.Where("name like ?", "%"+*query+"%")
	}

	queries.Order("name desc").Find(&productCategories)

	return productCategories
}

func (repository *BrandRepositoryImpl) GetById(db *gorm.DB, brandId uint) (entity.Brand, error) {
	id := brandId
	var brands entity.Brand
	err := db.
		Preload("CreatedBy").
		Preload("ModifiedBy").
		First(&brands, "id = ?", id).Error

	return brands, err
}
