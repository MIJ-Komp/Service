package response

import (
	"time"
)

type ProductCategoryResponse struct {
	Id           uint      `json:"id"`
	Name         string    `json:"name"`
	CreatedById  uint      `json:"createdById"`
	CreatedAt    time.Time `json:"createdAt"`
	ModifiedById uint      `json:"modifiedById"`
	ModifiedAt   time.Time `json:"modifiedAt"`
}
