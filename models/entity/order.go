package entity

import (
	"time"

	"api.mijkomp.com/models/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	Id         uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code       string            `gorm:"type:VARCHAR(30); null; unique;"`
	CustomerId *uint             `gorm:"type:bigint; null;"`
	OrderDate  time.Time         `gorm:"type:timestamptz; not null"`
	Status     enum.EOrderStatus `gorm:"type:varchar(10); not null;"`
	Notes      string            `gorm:"type:varchar(4096); null;"`
	IsPaid     bool              `gorm:"null;"`
	TotalPaid  *float64          `gorm:"type:decimal(17,5); null"`
	PaidAt     *time.Time        `gorm:"type:timestamptz; null"`
	PaymentId  *uuid.UUID        `gorm:"type:uuid; null;"`
	DeletedAt  gorm.DeletedAt    `gorm:"index"`

	CreatedByCustomerAt  time.Time `gorm:"type:timestamptz; null"`
	ModifiedByCustomerAt time.Time `gorm:"type:timestamptz; null"`

	ModifiedByAdminAt *time.Time `gorm:"type:timestamptz; null"`
	ModifiedByAdminId *uint      `gorm:"type:bigint; null; foreign_key;"`

	OrderItems   []OrderItem  `gorm:"foreignKey:order_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerInfo CustomerInfo `gorm:"foreignKey:order_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ShippingInfo ShippingInfo `gorm:"foreignKey:order_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Payment      Payment      `gorm:"foreignKey:PaymentId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ModifiedByAdmin User
}

type OrderItem struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderId      uuid.UUID `gorm:"foreignKey; type:uuid; not null;"`
	ProductId    uuid.UUID `gorm:"foreignKey; type:uuid; not null;"`
	ProductSkuId uuid.UUID `gorm:"foreignKey; type:uuid; not null;"`
	Quantity     int       `gorm:"type:int; not null;"`
	Price        float64   `gorm:"type:decimal(17,5); not null;"`
	Sequence     int       `gorm:"type:int; not null;"`
	CreatedAt    time.Time `gorm:"type:timestamptz; not null"`

	Product Product `gorm:"foreignKey:product_id;references:Id;"`
}

type CustomerInfo struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderId     uuid.UUID `gorm:"foreignKey; type:uuid; not null;"`
	Name        string    `gorm:"type:varchar(128); not null;"`
	Email       string    `gorm:"type:varchar(128); not null;"`
	PhoneNumber string    `gorm:"type:varchar(15); not null;"`
}

type ShippingInfo struct {
	Id                uuid.UUID  `gorm:"type:uuid; primaryKey; default:gen_random_uuid();"`
	OrderId           uuid.UUID  `gorm:"foreignKey; type:uuid; not null;"`
	RecipientName     string     `gorm:"type:varchar(256); not null"`
	Address           string     `gorm:"type:vaarchar(4086); not null"`
	District          string     `gorm:"type:varchar(128); not null"`
	City              string     `gorm:"type:varchar(128); not null"`
	Province          string     `gorm:"type:varchar(128); not null"`
	PostalCode        string     `gorm:"type:varchar(128); not null"`
	ShippingMethod    *string    `gorm:"type:varchar(256); null"`
	TrackingNumber    *string    `gorm:"type:varchar(256); null"`
	EstimatedDelivery *time.Time `gorm:"type:timestamptz; null"`
	ShippedAt         *time.Time `gorm:"type:timestamptz; null"`
	DeliveredAt       *time.Time `gorm:"type:timestamptz; null"`
	Notes             *string    `gorm:"type:varcha(2048); null"`
	CreatedAt         time.Time  `gorm:"type:timestamptz; not null"`
	UpdatedAt         time.Time  `gorm:"type:timestamptz; not null"`
}
