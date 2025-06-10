package repository_impl

import (
	"api.mijkomp.com/models/entity"
	"gorm.io/gorm"
)

type ProductCategoryRepositoryImpl struct {
}

func NewProductCategoryRepository() *ProductCategoryRepositoryImpl {
	return &ProductCategoryRepositoryImpl{}
}

func (repository *ProductCategoryRepositoryImpl) Save(db *gorm.DB, productCategory entity.ProductCategory) (entity.ProductCategory, error) {
	err := db.Save(&productCategory).Error
	return productCategory, err
}

func (repository *ProductCategoryRepositoryImpl) Delete(db *gorm.DB, productCategoryId uint) error {
	err := db.Where("id = ?", productCategoryId).Delete(&entity.ProductCategory{}).Error
	return err
}

func (repository *ProductCategoryRepositoryImpl) Search(db *gorm.DB, query *string, parentId *uint) []entity.ProductCategory {

	productCategories := []entity.ProductCategory{}

	queries := db.Model(&entity.ProductCategory{})

	if query != nil {
		queries.Where("name like ?", "%"+*query+"%")
	}

	if parentId != nil {
		queries.Where("parent_id = ?", parentId)
	}
	queries.Order("product_categories.name desc").Find(&productCategories)

	return productCategories
}

func (repository *ProductCategoryRepositoryImpl) GetById(db *gorm.DB, productCategoryId uint) (entity.ProductCategory, error) {
	id := productCategoryId
	var productCategorys entity.ProductCategory
	err := db.
		Preload("CreatedBy").
		Preload("ModifiedBy").
		First(&productCategorys, "id = ?", id).Error

	return productCategorys, err
}
