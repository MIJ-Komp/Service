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
	Product      Product   `gorm:"foreignKey:product_id;references:Id;"`
}

type CustomerInfo struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderId     uuid.UUID `gorm:"foreignKey; type:uuid; not null;"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
}

type ShippingInfo struct {
	Id                uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderId           uuid.UUID  `gorm:"foreignKey; type:uuid; not null;"`
	RecipientName     string     `gorm:"size:100;not null"`
	Address           string     `gorm:"type:text;not null"`
	Province          string     `gorm:"size:100"`
	City              string     `gorm:"size:100"`
	PostalCode        string     `gorm:"size:20"`
	ShippingMethod    *string    `gorm:"size:50"`
	TrackingNumber    *string    `gorm:"size:100"`
	EstimatedDelivery *time.Time `gorm:"type:date"`
	ShippedAt         *time.Time `gorm:"type:timestamptz"`
	DeliveredAt       *time.Time `gorm:"type:timestamptz"`
	Notes             *string    `gorm:"type:text"`
	CreatedAt         time.Time  `gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime"`
}
