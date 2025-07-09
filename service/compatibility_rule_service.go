package service

import (
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
)

type CompatibilityRuleService interface {
	Create(currentUserId uint, payload request.CompatibilityRule) response.CompatibilityRule
	Update(currentUserId uint, CompatibilityRuleId uint, payload request.CompatibilityRule) response.CompatibilityRule
	Delete(currentUserId uint, CompatibilityRuleId uint) string
	Search(currentUserId uint, sourceComponentTypeCode *string, targetComponentTypeCode *string) []response.CompatibilityRule
	GetById(currentUserId uint, CompatibilityRuleId uint) response.CompatibilityRule
}
