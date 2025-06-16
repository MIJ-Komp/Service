package response

import "time"

type Menu struct {
	Id         uint       `json:"id"`
	Name       string     `json:"name"`
	ParentId   *uint      `json:"parentId"`
	CreatedBy  AuditTrail `json:"createdById"`
	CreatedAt  time.Time  `json:"createdAt"`
	ModifiedBy AuditTrail `json:"modifiedById"`
	ModifiedAt time.Time  `json:"modifiedAt"`

	MenuItems []MenuItem `json:"menuItems"`
}

type MenuItem struct {
	Id                uint   `json:"id"`
	ProductCategoryId uint   `json:"productCategoryId"`
	Name              string `json:"name"`
}
