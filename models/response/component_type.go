package response

import (
	"time"
)

type ComponentType struct {
	Id          uint       `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedBy   AuditTrail `json:"createdById"`
	CreatedAt   time.Time  `json:"createdAt"`
	ModifiedBy  AuditTrail `json:"modifiedById"`
	ModifiedAt  time.Time  `json:"modifiedAt"`
}
