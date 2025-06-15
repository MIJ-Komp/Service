package response

import (
	"time"
)

type CompatibilityRule struct {
	Id                      uint   `json:"id"`
	SourceComponentTypeCode string `json:"sourceComponentTypeCode"`
	TargetComponentTypeCode string `json:"targetComponentTypeCode"`
	SourceKey               string `json:"sourceKey"`
	TargetKey               string `json:"targetKey"`
	Condition               string `json:"condition"`
	ValueType               string `json:"valueType"`
	ErrorMessage            string `json:"errorMessage"`
	IsActive                bool   `json:"isActive"`

	CreatedBy  AuditTrail `json:"createdById"`
	CreatedAt  time.Time  `json:"createdAt"`
	ModifiedBy AuditTrail `json:"modifiedById"`
	ModifiedAt time.Time  `json:"modifiedAt"`
}
