package repository

import (
	"api.mijkomp.com/models/entity"
	"gorm.io/gorm"
)

type CompatibilityRuleRepository interface {
	Save(db *gorm.DB, payload entity.CompatibilityRule) (entity.CompatibilityRule, error)
	Delete(db *gorm.DB, compatibilityRuleId uint) error
	Search(db *gorm.DB, query *string) []entity.CompatibilityRule
	GetById(db *gorm.DB, compatibilityRuleId uint) (entity.CompatibilityRule, error)
}
