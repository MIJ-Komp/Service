package response

import (
	"time"

	"github.com/google/uuid"
)

type ProductResponse struct {
	Id                      uuid.UUID    `json:"id"`
	ProductType             EnumResponse `json:"productType"`
	Name                    string       `json:"name"`
	SKU                     string       `json:"sku"`
	IsActive                bool         `json:"isActive"`
	IsShowOnlyInMarketPlace bool         `json:"isShowOnlyInMarketPlace"`
	ImageIds                []uuid.UUID  `json:"imageIds"`
	VideoUrl                *string      `json:"videoUrl"`
	Description             string       `json:"description"`
	ProductSpec             string       `json:"productSpec"`
	CreatedById             uint         `json:"createdById"`
	CreatedAt               time.Time    `json:"createdAt"`
	ModifiedById            uint         `json:"modifiedById"`
	ModifiedAt              time.Time    `json:"modifiedAt"`
	Tags                    *string      `json:"tags"`

	ProductCategory *ProductCategory `json:"productCategory"`
	Brand           *Brand           `json:"brand"`

	ProductSkus []ProductSku `json:"productSkus"`

	ProductGroupItems []ProductGroupItemResponse `json:"productGroupItems"`

	ProductVariantOptions      []ProductVariantOption      `json:"productVariantOptions"`
	ProductVariantOptionValues []ProductVariantOptionValue `json:"productVariantOptionValues"`

	ProductSkuVariants []ProductSkuVariant `json:"productSkuVariants"`
}

type ProductSku struct {
	Id              uuid.UUID `json:"id"`
	ProductId       uuid.UUID `json:"productId"`
	SKU             string    `json:"sku"`
	Name            string    `json:"name"`
	Price           float64   `json:"price"`
	Stock           *int      `json:"stock"`
	StockAlert      *int      `json:"stockAlert"`
	Sequence        int       `json:"sequence"`
	ComponentTypeId *uint     `json:"componentTypeId"`
	IsActive        bool      `json:"isActive"`

	ComponentSpecs    []ComponentSpec            `json:"componentSpecs"`
	ProductGroupItems []ProductGroupItemResponse `json:"productGroupItems"`
}

type ComponentSpec struct {
	Id           uuid.UUID `json:"id"`
	ProductSkuId uuid.UUID `json:"productSkuId"`
	SpecKey      string    `json:"specKey"`
	SpecValue    string    `json:"specValue"`
}

type ProductVariantOption struct {
	Id          uuid.UUID `json:"id"`
	ProductId   uuid.UUID `json:"productId"`
	Name        string    `json:"name"`
	AllowDelete bool      `json:"allowDelete"`
	Sequence    int       `json:"sequence"`
}

type ProductVariantOptionValue struct {
	Id                     uuid.UUID `json:"id"`
	ProductVariantOptionId uuid.UUID `json:"productVariantOptionId"`
	Name                   string    `json:"name"`
	AllowDelete            bool      `json:"allowDelete"`
	Sequence               int       `json:"sequence"`
}

type ProductSkuVariant struct {
	Id                          uuid.UUID `json:"id"`
	ProductSkuId                uuid.UUID `json:"productSkuId"`
	ProductVariantOptionId      uuid.UUID `json:"productVariantOptionId"`
	ProductVariantOptionValueId uuid.UUID `json:"productVariantOptionValueId"`
}

type ProductGroupItemResponse struct {
	Id           uuid.UUID       `json:"id"`
	ProductId    uuid.UUID       `json:"productId"`
	ProductSkuId uuid.UUID       `json:"productSkuId"`
	Qty          int             `json:"qty"`
	Product      ProductResponse `json:"product"`
}

type VariantOption struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type BrowseProductSku struct {
	Id           uuid.UUID  `json:"id"`
	ProductId    uuid.UUID  `json:"productId"`
	SKU          string     `json:"sku"`
	Name         string     `json:"name"`
	Stock        *int       `json:"stock"`
	Price        float64    `json:"price"`
	IsActive     bool       `json:"isActive"`
	ProductType  string     `json:"productType"`
	PictureId    *uuid.UUID `json:"pictureId"`
	Description  string     `json:"description"`
	CreatedById  uint       `json:"createdById"`
	CreatedAt    time.Time  `json:"createdAt"`
	ModifiedById uint       `json:"modifiedById"`
	ModifiedAt   time.Time  `json:"modifiedAt"`
}
