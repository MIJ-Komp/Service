package repository_impl

import (
	"api.mijkomp.com/models/entity"
	"gorm.io/gorm"
)

type CompatibilityRuleRepositoryImpl struct {
}

func NewCompatibilityRuleRepository() *CompatibilityRuleRepositoryImpl {
	return &CompatibilityRuleRepositoryImpl{}
}

func (repository *CompatibilityRuleRepositoryImpl) Save(db *gorm.DB, compatibilityRule entity.CompatibilityRule) (entity.CompatibilityRule, error) {
	err := db.Save(&compatibilityRule).Error
	return compatibilityRule, err
}

func (repository *CompatibilityRuleRepositoryImpl) Delete(db *gorm.DB, compatibilityRuleId uint) error {
	err := db.Where("id = ?", compatibilityRuleId).Delete(&entity.CompatibilityRule{}).Error
	return err
}

func (repository *CompatibilityRuleRepositoryImpl) Search(db *gorm.DB, sourceComponentTypeCode *string, targetComponentTypeCode *string) []entity.CompatibilityRule {

	compatibilityRules := []entity.CompatibilityRule{}

	queries := db.Model(&entity.CompatibilityRule{})

	if sourceComponentTypeCode != nil {
		queries.Where("name source_component_type_code = ?", sourceComponentTypeCode)
	}

	if targetComponentTypeCode != nil {
		queries.Where("name target_component_type_code = ?", targetComponentTypeCode)
	}

	queries.Order("modified_at desc").Find(&compatibilityRules)

	return compatibilityRules
}

func (repository *CompatibilityRuleRepositoryImpl) GetById(db *gorm.DB, compatibilityRuleId uint) (entity.CompatibilityRule, error) {
	id := compatibilityRuleId
	var compatibilityRules entity.CompatibilityRule
	err := db.
		Preload("CreatedBy").
		Preload("ModifiedBy").
		First(&compatibilityRules, "id = ?", id).Error

	return compatibilityRules, err
}
