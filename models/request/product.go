package request

import (
	"api.mijkomp.com/models/enum"
	"github.com/google/uuid"
)

type ProductPayload struct {
	ProductType             enum.EProductType `json:"productType"`
	SKU                     string            `json:"sku"`
	Name                    string            `json:"name"`
	IsActive                bool              `json:"isActive"`
	IsShowOnlyInMarketPlace bool              `json:"isShowOnlyInMarketPlace"`
	ImageIds                []uuid.UUID       `json:"imageIds"`
	ProductCategoryId       *uint             `json:"productCategoryId"`
	Description             string            `json:"description"`

	ProductSkus        []ProductSkuPayload        `json:"productSkus"`
	ProductSkuVariants []ProductSkuVariantPayload `json:"productSkuVariants"`

	ProductVariantOptions      []ProductVariantOptionPayload      `json:"productVariantOptions"`
	ProductVariantOptionValues []ProductVariantOptionValuePayload `json:"productVariantOptionValues"`

	ProductGroupItems []ProductGroupItemPayload `json:"productGroupItems"`
}

type ProductVariantOptionPayload struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ProductVariantOptionValuePayload struct {
	Id                     uuid.UUID `json:"id"`
	ProductVariantOptionId uuid.UUID `json:"productVariantOptionId"`
	Name                   string    `json:"name"`
}

type ProductSkuPayload struct {
	Id           uuid.UUID     `json:"id"`
	SKU          string        `json:"sku"`
	Name         string        `json:"name"`
	Price        float64       `json:"price"`
	Stock        *int          `json:"stock"`
	StockAlert   *int          `json:"stockAlert"`
	IsActive     bool          `json:"isActive"`
	ProductSpecs []ProductSpec `json:"productSpecs"`
}

type ProductSpec struct {
	Id        uuid.UUID `json:"id"`
	SpecKey   string    `json:"specKey"`
	SpecValue string    `json:"specValue"`
}

type ProductSkuVariantPayload struct {
	Id                          uuid.UUID `json:"id"`
	ProductSkuId                uuid.UUID `json:"productSkuId"`
	ProductVariantOptionId      uuid.UUID `json:"productVariantOptionId"`
	ProductVariantOptionValueId uuid.UUID `json:"productVariantOptionValueId"`
}

type ProductGroupItemPayload struct {
	Id           uuid.UUID `json:"id"`
	ProductId    uuid.UUID `json:"productId"`
	ProductSkuId uuid.UUID `json:"productSkuId"`
	Qty          int       `json:"qty"`
}

type VariantOptionPayload struct {
	Name string `json:"name"`
}
