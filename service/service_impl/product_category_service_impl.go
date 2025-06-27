package service_impl

import (
	"time"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers"
	"api.mijkomp.com/helpers/data_mapper"
	"api.mijkomp.com/models/entity"
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/repository"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ProductCategoryServiceImpl struct {
	ProductCategoryRepository repository.ProductCategoryRepository
	db                        *gorm.DB
	Validation                *validator.Validate
}

func NewProductCategoryService(categoryRepostitory repository.ProductCategoryRepository, validation *validator.Validate, db *gorm.DB) *ProductCategoryServiceImpl {
	return &ProductCategoryServiceImpl{
		ProductCategoryRepository: categoryRepostitory,
		Validation:                validation,
		db:                        db,
	}
}

func (service *ProductCategoryServiceImpl) Create(currentUserId uint, payload request.ProductCategory) response.ProductCategory {

	err := service.Validation.Struct(payload)
	exception.PanicIfNeeded(err)

	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	categoryEntity := entity.ProductCategory{
		Name:         payload.Name,
		CreatedById:  currentUserId,
		CreatedAt:    time.Now(),
		ModifiedById: currentUserId,
		ModifiedAt:   time.Now(),
	}

	result, err := service.ProductCategoryRepository.Save(tx, categoryEntity)
	exception.PanicIfNeeded(err)

	return service.GenerateGetResult(result)
}

func (service *ProductCategoryServiceImpl) Update(currentUserId uint, categoryId uint, payload request.ProductCategory) response.ProductCategory {

	err := service.Validation.Struct(payload)
	exception.PanicIfNeeded(err)

	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	category, err := service.ProductCategoryRepository.GetById(tx, categoryId)
	exception.PanicIfNeeded(err)

	category.Name = payload.Name
	category.ModifiedById = currentUserId
	category.ModifiedAt = time.Now()

	result, err := service.ProductCategoryRepository.Save(tx, category)
	exception.PanicIfNeeded(err)

	return service.GenerateGetResult(result)
}

func (service *ProductCategoryServiceImpl) Delete(currentUserId uint, categoryId uint) string {
	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	err := service.ProductCategoryRepository.Delete(tx, categoryId)
	exception.PanicIfNeeded(err)

	return "Kategori berhasil dihapus"
}

func (service *ProductCategoryServiceImpl) Search(currentUserId uint, query *string) []response.ProductCategory {
	res := service.ProductCategoryRepository.Search(service.db, query)

	return service.GenerateSearchResult(res)
}

func (service *ProductCategoryServiceImpl) GetById(currentUserId uint, categoryId uint) response.ProductCategory {
	res, err := service.ProductCategoryRepository.GetById(service.db, categoryId)
	exception.PanicIfNeeded(err)

	return service.GenerateGetResult(res)
}

// map helpers

func (service *ProductCategoryServiceImpl) GenerateSearchResult(categories []entity.ProductCategory) []response.ProductCategory {

	res := []response.ProductCategory{}
	for _, category := range categories {
		res = append(res, service.GenerateGetResult(category))
	}

	return res
}

func (service *ProductCategoryServiceImpl) GenerateGetResult(category entity.ProductCategory) response.ProductCategory {
	res := response.ProductCategory{
		Id:         category.Id,
		Name:       category.Name,
		CreatedBy:  data_mapper.MapAuditTrail(category.CreatedBy),
		CreatedAt:  category.CreatedAt,
		ModifiedBy: data_mapper.MapAuditTrail(category.ModifiedBy),
		ModifiedAt: category.ModifiedAt,
	}

	return res
}
