package entity

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	Id            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	IsPaid        bool       `gorm:"not null"`
	TotalPaid     *float64   `gorm:"type:decimal(17,5);"`
	PaidAt        *time.Time `gorm:"type:timestamptz; null;"`
	InvoiceUrl    *string    `gorm:"type:varchar(128)"`
	InvoiceId     *string    `gorm:"type:varchar(128)"`
	PayerInfo     *string    `gorm:"type:varchar(128)"`
	PaymentMethod *string    `gorm:"type:varchar(128)"`
	// CreatedById   uuid.UUID `gorm:"type:uniqueidentifier; foreignKey; not null;"`
	// CreatedAt     time.Time              `gorm:"type:datetime; not null;"`
	// ModifiedById  uuid.UUID `gorm:"type:uniqueidentifier; foreignKey; not null;"`
	// ModifiedAt    time.Time              `gorm:"type:datetime; not null;"`
}
