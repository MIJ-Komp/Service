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
	IsPartOfCompositeOnly   bool         `json:"isPartOfCompositeOnly"`
	IsTrackInventory        bool         `json:"isTrackInventory"`
	IsPriceSameForAllOutlet bool         `json:"isPriceSameForAllOutlet"`
	PictureId               *uuid.UUID   `json:"pictureId"`
	Description             string       `json:"description"`
	CreatedById             uint         `json:"createdById"`
	CreatedAt               time.Time    `json:"createdAt"`
	ModifiedById            uint         `json:"modifiedById"`
	ModifiedAt              time.Time    `json:"modifiedAt"`

	ProductCategory *ProductCategoryResponse `json:"productCategory"`

	ProductSkus []ProductSku `json:"productSkus"`

	ProductGroupItems []ProductGroupItemResponse `json:"productGroupItems"`

	ProductVariantOptions      []ProductVariantOption      `json:"productVariantOptions"`
	ProductVariantOptionValues []ProductVariantOptionValue `json:"productVariantOptionValues"`

	ProductSkuVariants []ProductSkuVariant `json:"productSkuVariants"`
}

type ProductSku struct {
	Id        uuid.UUID `json:"id"`
	ProductId uuid.UUID `json:"productId"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	Sequence  int       `json:"sequence"`
	IsActive  bool      `json:"isActive"`
	// ProductSkuDetails []ProductSkuDetail `json:"productSkuDetails"`
}

// type ProductSkuDetail struct {
// 	Id           uint     `json:"id"`
// 	ProductSkuId uint     `json:"productSkuId"`
// 	OutletId     uint     `json:"outletId"`
// 	Fee          *float64 `json:"fee"`
// 	CostPrice    *float64 `json:"costPrice"`
// 	SellingPrice float64  `json:"sellingPrice"`
// 	Qty          *int     `json:"qty"`
// 	QtyAlert     *int     `json:"qtyAlert"`
// }

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
