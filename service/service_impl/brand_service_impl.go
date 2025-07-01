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

type BrandServiceImpl struct {
	BrandRepository repository.BrandRepository
	db              *gorm.DB
	Validation      *validator.Validate
}

func NewBrandService(brandRepostitory repository.BrandRepository, validation *validator.Validate, db *gorm.DB) *BrandServiceImpl {
	return &BrandServiceImpl{
		BrandRepository: brandRepostitory,
		Validation:      validation,
		db:              db,
	}
}

func (service *BrandServiceImpl) Create(currentUserId uint, payload request.Brand) response.Brand {

	err := service.Validation.Struct(payload)
	exception.PanicIfNeeded(err)

	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	brandEntity := entity.Brand{
		Name:         payload.Name,
		CreatedById:  currentUserId,
		CreatedAt:    time.Now(),
		ModifiedById: currentUserId,
		ModifiedAt:   time.Now(),
	}

	result, err := service.BrandRepository.Save(tx, brandEntity)
	exception.PanicIfNeeded(err)

	return service.GenerateGetResult(result)
}

func (service *BrandServiceImpl) Update(currentUserId uint, brandId uint, payload request.Brand) response.Brand {

	err := service.Validation.Struct(payload)
	exception.PanicIfNeeded(err)

	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	brand, err := service.BrandRepository.GetById(tx, brandId)
	exception.PanicIfNeeded(err)

	brand.Name = payload.Name
	brand.ModifiedById = currentUserId
	brand.ModifiedAt = time.Now()

	result, err := service.BrandRepository.Save(tx, brand)
	exception.PanicIfNeeded(err)

	return service.GenerateGetResult(result)
}

func (service *BrandServiceImpl) Delete(currentUserId uint, brandId uint) string {
	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	err := service.BrandRepository.Delete(tx, brandId)
	exception.PanicIfNeeded(err)

	return "Kategori berhasil dihapus"
}

func (service *BrandServiceImpl) Search(currentUserId uint, query *string) []response.Brand {
	res := service.BrandRepository.Search(service.db, query)

	return service.GenerateSearchResult(res)
}

func (service *BrandServiceImpl) GetById(currentUserId uint, brandId uint) response.Brand {
	res, err := service.BrandRepository.GetById(service.db, brandId)
	exception.PanicIfNeeded(err)

	return service.GenerateGetResult(res)
}

// map helpers

func (service *BrandServiceImpl) GenerateSearchResult(categories []entity.Brand) []response.Brand {

	res := []response.Brand{}
	for _, brand := range categories {
		res = append(res, service.GenerateGetResult(brand))
	}

	return res
}

func (service *BrandServiceImpl) GenerateGetResult(brand entity.Brand) response.Brand {
	res := response.Brand{
		Id:         brand.Id,
		Name:       brand.Name,
		CreatedBy:  data_mapper.MapAuditTrail(brand.CreatedBy),
		CreatedAt:  brand.CreatedAt,
		ModifiedBy: data_mapper.MapAuditTrail(brand.ModifiedBy),
		ModifiedAt: brand.ModifiedAt,
	}

	return res
}
