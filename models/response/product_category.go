package response

import "time"

type ProductCategory struct {
	Id         uint       `json:"id"`
	Name       string     `json:"name"`
	ParentId   *uint      `json:"parentId"`
	CreatedBy  AuditTrail `json:"createdById"`
	CreatedAt  time.Time  `json:"createdAt"`
	ModifiedBy AuditTrail `json:"modifiedById"`
	ModifiedAt time.Time  `json:"modifiedAt"`
}
