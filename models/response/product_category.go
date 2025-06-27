package response

import "time"

type ProductCategory struct {
	Id         uint       `json:"id"`
	Name       string     `json:"name"`
	CreatedBy  AuditTrail `json:"createdById"`
	CreatedAt  time.Time  `json:"createdAt"`
	ModifiedBy AuditTrail `json:"modifiedById"`
	ModifiedAt time.Time  `json:"modifiedAt"`
}

type Brand struct {
	Id         uint       `json:"id"`
	Name       string     `json:"name"`
	CreatedBy  AuditTrail `json:"createdById"`
	CreatedAt  time.Time  `json:"createdAt"`
	ModifiedBy AuditTrail `json:"modifiedById"`
	ModifiedAt time.Time  `json:"modifiedAt"`
}
