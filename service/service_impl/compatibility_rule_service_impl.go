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

type CompatibilityRuleServiceImpl struct {
	CompatibilityRuleRepository repository.CompatibilityRuleRepository
	db                          *gorm.DB
	Validation                  *validator.Validate
}

func NewCompatibilityRuleService(compatibilityRuleRepostitory repository.CompatibilityRuleRepository, validation *validator.Validate, db *gorm.DB) *CompatibilityRuleServiceImpl {
	return &CompatibilityRuleServiceImpl{
		CompatibilityRuleRepository: compatibilityRuleRepostitory,
		Validation:                  validation,
		db:                          db,
	}
}

func (service *CompatibilityRuleServiceImpl) Create(currentUserId uint, payload request.CompatibilityRule) response.CompatibilityRule {

	err := service.Validation.Struct(payload)
	exception.PanicIfNeeded(err)

	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	compatibilityRuleEntity := entity.CompatibilityRule{
		SourceComponentTypeCode: payload.SourceComponentTypeCode,
		TargetComponentTypeCode: payload.TargetComponentTypeCode,
		SourceKey:               payload.SourceKey,
		TargetKey:               payload.TargetKey,
		Condition:               payload.Condition,
		ValueType:               payload.ValueType,
		ErrorMessage:            payload.ErrorMessage,
		IsActive:                payload.IsActive,
		CreatedById:             currentUserId,
		CreatedAt:               time.Now(),
		ModifiedById:            currentUserId,
		ModifiedAt:              time.Now(),
	}

	result, err := service.CompatibilityRuleRepository.Save(tx, compatibilityRuleEntity)
	exception.PanicIfNeeded(err)

	return service.GenerateGetResult(result)
}

func (service *CompatibilityRuleServiceImpl) Update(currentUserId uint, compatibilityRuleId uint, payload request.CompatibilityRule) response.CompatibilityRule {

	err := service.Validation.Struct(payload)
	exception.PanicIfNeeded(err)

	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	compatibilityRule, err := service.CompatibilityRuleRepository.GetById(tx, compatibilityRuleId)
	exception.PanicIfNeeded(err)
	compatibilityRule.SourceComponentTypeCode = payload.SourceComponentTypeCode
	compatibilityRule.TargetComponentTypeCode = payload.TargetComponentTypeCode
	compatibilityRule.SourceKey = payload.SourceKey
	compatibilityRule.TargetKey = payload.TargetKey
	compatibilityRule.Condition = payload.Condition
	compatibilityRule.ValueType = payload.ValueType
	compatibilityRule.ErrorMessage = payload.ErrorMessage
	compatibilityRule.IsActive = payload.IsActive

	compatibilityRule.ModifiedById = currentUserId
	compatibilityRule.ModifiedAt = time.Now()

	result, err := service.CompatibilityRuleRepository.Save(tx, compatibilityRule)
	exception.PanicIfNeeded(err)

	return service.GenerateGetResult(result)
}

func (service *CompatibilityRuleServiceImpl) Delete(currentUserId uint, compatibilityRuleId uint) string {
	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	err := service.CompatibilityRuleRepository.Delete(tx, compatibilityRuleId)
	exception.PanicIfNeeded(err)

	return "Kategori berhasil dihapus"
}

func (service *CompatibilityRuleServiceImpl) Search(currentUserId uint, query *string) []response.CompatibilityRule {
	res := service.CompatibilityRuleRepository.Search(service.db, query)

	return service.GenerateSearchResult(res)
}

func (service *CompatibilityRuleServiceImpl) GetById(currentUserId uint, compatibilityRuleId uint) response.CompatibilityRule {
	res, err := service.CompatibilityRuleRepository.GetById(service.db, compatibilityRuleId)
	exception.PanicIfNeeded(err)

	return service.GenerateGetResult(res)
}

// map helpers

func (service *CompatibilityRuleServiceImpl) GenerateSearchResult(compatibilityRules []entity.CompatibilityRule) []response.CompatibilityRule {

	res := []response.CompatibilityRule{}
	for _, compatibilityRule := range compatibilityRules {
		res = append(res, service.GenerateGetResult(compatibilityRule))
	}

	return res
}

func (service *CompatibilityRuleServiceImpl) GenerateGetResult(compatibilityRule entity.CompatibilityRule) response.CompatibilityRule {
	res := response.CompatibilityRule{
		Id:                      compatibilityRule.Id,
		SourceComponentTypeCode: compatibilityRule.SourceComponentTypeCode,
		TargetComponentTypeCode: compatibilityRule.TargetComponentTypeCode,
		SourceKey:               compatibilityRule.SourceKey,
		TargetKey:               compatibilityRule.TargetKey,
		Condition:               compatibilityRule.Condition,
		ValueType:               compatibilityRule.ValueType,
		ErrorMessage:            compatibilityRule.ErrorMessage,
		IsActive:                compatibilityRule.IsActive,
		CreatedBy:               data_mapper.MapAuditTrail(compatibilityRule.CreatedBy),
		CreatedAt:               compatibilityRule.CreatedAt,
		ModifiedBy:              data_mapper.MapAuditTrail(compatibilityRule.ModifiedBy),
		ModifiedAt:              compatibilityRule.ModifiedAt,
	}

	return res
}
