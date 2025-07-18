package repository

import (
	"api.mijkomp.com/models/entity"
	"api.mijkomp.com/models/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Save(db *gorm.DB, product entity.Product) (entity.Product, error)
	Delete(db *gorm.DB, product entity.Product) error
	Search(
		db *gorm.DB,
		isAdmin bool,
		query *string,
		ids *[]uuid.UUID,
		productTypes *[]string,
		componentTypeIds *[]uint,
		productCategoryIds *[]uint,
		brandIds *[]uint,
		tag *string,
		minPrice, maxPrice *float64,
		isInStockOnly *bool,
		isActive, isShowOnlyInMarketPlace *bool,
		page, pageSize *int,
	) ([]entity.Product, int64, int64)

	GetById(db *gorm.DB, productId uuid.UUID) (entity.Product, error)

	// Product SKU
	SaveProductSkus(db *gorm.DB, variant []entity.ProductSku) error
	DeleteProductSkus(db *gorm.DB, productId uuid.UUID, variants []entity.ProductSku) error
	UpdateStockProductSkus(db *gorm.DB, productSkus []entity.ProductSku) error
	BrowseProductSku(db *gorm.DB, query *string, productTypes *[]string, productCategoryId *uint, page, pageSize *int) ([]response.BrowseProductSku, int64, int64)
	GetProductSkuByIds(db *gorm.DB, productSkuIds []uuid.UUID) []response.BrowseProductSku
	// GetProductSkuByIds(db *gorm.DB, productSkuIds []uuid.UUID) []entity.ProductSku

	// Product Component Specs
	SaveComponentSpecs(db *gorm.DB, componentSpecs []entity.ComponentSpec) error
	DeleteComponentSpecs(db *gorm.DB, componentSpecs []entity.ComponentSpec) error

	// Group Details
	SaveProductGroupItems(db *gorm.DB, groupItems []entity.ProductGroupItem) error
	DeleteProductGroupItems(db *gorm.DB, groupItems []entity.ProductGroupItem) error

	// Variant Options
	SaveProductVariantOptions(db *gorm.DB, variantOptions []entity.ProductVariantOption) error
	DeleteProductVariantOptions(db *gorm.DB, variantOptions []entity.ProductVariantOption) error

	// Variant Option Values
	SaveProductVariantOptionValues(db *gorm.DB, optionValues []entity.ProductVariantOptionValue) error
	DeleteProductVariantOptionValues(db *gorm.DB, optionValues []entity.ProductVariantOptionValue) error

	// Product SKU Variant
	SaveProductSkuVariants(db *gorm.DB, productSkuVariants []entity.ProductSkuVariant) error
	DeleteProductSkuVariants(db *gorm.DB, productSkuVariants []entity.ProductSkuVariant) error

	// Variant Option
	SaveVariantOption(db *gorm.DB, variantOption entity.VariantOption) (entity.VariantOption, error)
	SaveVariantOptions(db *gorm.DB, variantOptions []entity.VariantOption) error
	GetVariantOptions(db *gorm.DB) ([]entity.VariantOption, error)
}
