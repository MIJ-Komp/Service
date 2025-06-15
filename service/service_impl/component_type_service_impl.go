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

type ComponentTypeServiceImpl struct {
	ComponentTypeRepository repository.ComponentTypeRepository
	db                      *gorm.DB
	Validation              *validator.Validate
}

func NewComponentTypeService(componentTypeRepostitory repository.ComponentTypeRepository, validation *validator.Validate, db *gorm.DB) *ComponentTypeServiceImpl {
	return &ComponentTypeServiceImpl{
		ComponentTypeRepository: componentTypeRepostitory,
		Validation:              validation,
		db:                      db,
	}
}

func (service *ComponentTypeServiceImpl) Create(currentUserId uint, payload request.ComponentType) response.ComponentType {

	err := service.Validation.Struct(payload)
	exception.PanicIfNeeded(err)

	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	componentTypeEntity := entity.ComponentType{
		Code:         payload.Code,
		Name:         payload.Name,
		Description:  payload.Description,
		CreatedById:  currentUserId,
		CreatedAt:    time.Now(),
		ModifiedById: currentUserId,
		ModifiedAt:   time.Now(),
	}

	result, err := service.ComponentTypeRepository.Save(tx, componentTypeEntity)
	exception.PanicIfNeeded(err)

	return service.GenerateGetResult(result)
}

func (service *ComponentTypeServiceImpl) Update(currentUserId uint, componentTypeId uint, payload request.ComponentType) response.ComponentType {

	err := service.Validation.Struct(payload)
	exception.PanicIfNeeded(err)

	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	componentType, err := service.ComponentTypeRepository.GetById(tx, componentTypeId)
	exception.PanicIfNeeded(err)

	componentType.Code = payload.Code
	componentType.Name = payload.Name
	componentType.Description = payload.Description
	componentType.ModifiedById = currentUserId
	componentType.ModifiedAt = time.Now()

	result, err := service.ComponentTypeRepository.Save(tx, componentType)
	exception.PanicIfNeeded(err)

	return service.GenerateGetResult(result)
}

func (service *ComponentTypeServiceImpl) Delete(currentUserId uint, componentTypeId uint) string {
	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	err := service.ComponentTypeRepository.Delete(tx, componentTypeId)
	exception.PanicIfNeeded(err)

	return "Tipe komponen berhasil dihapus"
}

func (service *ComponentTypeServiceImpl) Search(currentUserId uint, query *string) []response.ComponentType {
	res := service.ComponentTypeRepository.Search(service.db, query)

	return service.GenerateSearchResult(res)
}

func (service *ComponentTypeServiceImpl) GetById(currentUserId uint, componentTypeId uint) response.ComponentType {
	res, err := service.ComponentTypeRepository.GetById(service.db, componentTypeId)
	exception.PanicIfNeeded(err)

	return service.GenerateGetResult(res)
}

// map helpers

func (service *ComponentTypeServiceImpl) GenerateSearchResult(componentTypes []entity.ComponentType) []response.ComponentType {

	res := []response.ComponentType{}
	for _, componentType := range componentTypes {
		res = append(res, service.GenerateGetResult(componentType))
	}

	return res
}

func (service *ComponentTypeServiceImpl) GenerateGetResult(componentType entity.ComponentType) response.ComponentType {
	res := response.ComponentType{
		Id:          componentType.Id,
		Code:        componentType.Code,
		Name:        componentType.Name,
		Description: componentType.Description,
		CreatedBy:   data_mapper.MapAuditTrail(componentType.CreatedBy),
		CreatedAt:   componentType.CreatedAt,
		ModifiedBy:  data_mapper.MapAuditTrail(componentType.ModifiedBy),
		ModifiedAt:  componentType.ModifiedAt,
	}

	return res
}
