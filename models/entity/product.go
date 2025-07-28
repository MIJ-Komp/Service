package entity

import (
	"time"

	"api.mijkomp.com/models/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	Id                      uuid.UUID         `gorm:"primaryKey;type:uuid;"`
	ProductType             enum.EProductType `gorm:"type:varchar(8); not null;"`
	SKU                     string            `json:"type:varchar(256)"`
	Name                    string            `gorm:"type:varchar(256); not null;"`
	IsActive                bool              `gorm:"not null;"`
	IsShowOnlyInMarketPlace bool              `gorm:"not null;"`
	ImageIds                *string           `gorm:"type:varchar(2048); null;"`
	VideoUrl                *string           `gorm:"type:varchar(2048); null;"`
	ComponentTypeId         *uint             `gorm:"foreignKey; type:bigint; null;"`
	ProductCategoryId       *uint             `gorm:"type:bigint; foreignKey; null;"`
	BrandId                 *uint             `gorm:"type:bigint; foreignKey; null;"`
	Tags                    *string           `gorm:"type:varchar(256); null;"`
	Description             string            `gorm:"type:text; null;"`
	ProductSpec             string            `gorm:"type:varchar; null;"`
	CreatedById             uint              `gorm:"type:bigint; not null;"`
	CreatedAt               time.Time         `gorm:"type:timestamptz; not null;"`
	ModifiedById            uint              `gorm:"type:bigint; not null;"`
	ModifiedAt              time.Time         `gorm:"type:timestamptz; not null;"`
	DeletedAt               gorm.DeletedAt    `gorm:"index"`

	ProductCategory *ProductCategory
	Brand           *Brand
	ProductSkus     []ProductSku `gorm:"foreignKey:ProductId; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ComponentType   *ComponentType

	ProductVariantOptions      []ProductVariantOption      `gorm:"foreignKey:ProductId; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductVariantOptionValues []ProductVariantOptionValue `gorm:"foreignKey:ProductId; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductSkuVariants         []ProductSkuVariant         `gorm:"foreignKey:ProductId; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedBy  User
	ModifiedBy User
}

type ProductSku struct {
	Id         uuid.UUID      `gorm:"primaryKey; type:uuid;"`
	ProductId  uuid.UUID      `gorm:"foreignKey; type:uuid; not null;"`
	SKU        string         `gorm:"type:varchar(256); not null;"`
	Name       string         `gorm:"type: varchar(128); null"`
	Price      float64        `gorm:"type:decimal(17,5); not null;"`
	ImageId    *uuid.UUID     `gorm:"type:uuid; null;"`
	Stock      *int           `gorm:"type:integer; null"`
	StockAlert *int           `gorm:"type:integer; null"`
	IsActive   bool           `gorm:"not null;"`
	Sequence   int            `gorm:"type:int; not null"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	ComponentSpecs    []ComponentSpec    `gorm:"foreignKey:ProductSkuId; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductGroupItems []ProductGroupItem `gorm:"foreignKey:ParentId; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ComponentSpec struct {
	Id           uuid.UUID `gorm:"primaryKey; type:uuid;"`
	ProductSkuId uuid.UUID `gorm:"foreignKey; type:uuid; not null"`
	SpecKey      string    `gorm:"type:varchar(128); not null;"`
	SpecValue    string    `gorm:"type:varchar(128); not null;"`
	Sequence     int       `gorm:"type:int; not null"`
}

type ProductVariantOption struct {
	Id        uuid.UUID `gorm:"primaryKey; type:uuid;"`
	ProductId uuid.UUID `gorm:"foreignKey; type:uuid; not null"`
	Name      string    `gorm:"type:varchar(64);not null"`
	Sequence  int       `gorm:"type:integer; not null"`
}

type ProductVariantOptionValue struct {
	Id                     uuid.UUID `gorm:"primaryKey; type:uuid;"`
	ProductId              uuid.UUID `gorm:"foreignKey; type:uuid; not null"`
	ProductVariantOptionId uuid.UUID `gorm:"type:uuid; not null"`
	Name                   string    `gorm:"type:varchar(64); not null"`
	Sequence               int       `gorm:"type:integer; not null"`
}

type ProductSkuVariant struct {
	Id                          uuid.UUID `gorm:"primaryKey; type:uuid;"`
	ProductId                   uuid.UUID `gorm:"foreignKey; type:uuid; not null"`
	ProductSkuId                uuid.UUID `gorm:"type:uuid; not null"`
	ProductVariantOptionId      uuid.UUID `gorm:"type:uuid; not null"`
	ProductVariantOptionValueId uuid.UUID `gorm:"type:uuid; not null"`
	Sequence                    int       `gorm:"type:integer; not null"`
}

type ProductGroupItem struct {
	Id           uuid.UUID `gorm:"primaryKey; type:uuid;"`
	ParentId     uuid.UUID `gorm:"foreignKey:ParentId; type:uuid; not null"`
	ProductId    uuid.UUID `gorm:"foreignKey; type:uuid; not null"`
	ProductSkuId uuid.UUID `gorm:"foreignKey:ProductSkuId; type:uuid; not null"`
	Qty          int       `gorm:"type:integer; null"`
	Sequence     int       `gorm:"type:integer; not null"`

	Product Product `gorm:"foreignKey:ProductId;references:Id;"`
}

type VariantOption struct { // Master Data
	Id   uint   `gorm:"type:bigint; primaryKey; autoincrement;"`
	Name string `gorm:"type:varchar(64); not null"`
}
