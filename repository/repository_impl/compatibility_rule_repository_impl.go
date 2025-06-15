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

func (repository *CompatibilityRuleRepositoryImpl) Search(db *gorm.DB, query *string) []entity.CompatibilityRule {

	compatibilityRules := []entity.CompatibilityRule{}

	queries := db.Model(&entity.CompatibilityRule{})

	// if query != nil {
	// 	queries.Where("name like ?", "%"+*query+"%")
	// }

	// if parentId != nil {
	// 	queries.Where("parent_id = ?", parentId)
	// }
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
