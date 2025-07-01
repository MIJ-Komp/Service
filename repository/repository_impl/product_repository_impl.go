package repository_impl

import (
	"fmt"
	"strings"
	"time"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/models/entity"
	"api.mijkomp.com/models/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
}

func NewProductRepository() *ProductRepositoryImpl {
	return &ProductRepositoryImpl{}
}

// Product
func (repository *ProductRepositoryImpl) Save(db *gorm.DB, product entity.Product) (entity.Product, error) {
	var productOmitFields = []string{
		"ProductCategory",
		"Brand",
		"ProductSkus",
		"ProductVariantOptions",
		"ProductVariantOptionValues",
		"ProductSkuVariants",
		"ProductGroupItems",
	}
	err := db.Omit(
		productOmitFields...,
	).Save(&product).Error
	return product, err
}

func (repository *ProductRepositoryImpl) Delete(db *gorm.DB, product entity.Product) error {
	err := db.Where("id = ?", product.Id).Update("deleted_at", time.Now().UTC()).Error
	return err
}

func (repository *ProductRepositoryImpl) Search(
	db *gorm.DB,
	query *string,
	productTypes *[]string,
	productCategoryIds *[]uint,
	isActive, isShowOnlyInMarketPlace *bool,
	page, pageSize *int,
) ([]entity.Product, int64, int64) {
	var products []entity.Product
	var totalCount int64 = 0
	var totalPage int64 = 0

	queries := db.Model(&entity.Product{}).
		Preload("ProductCategory").
		Preload("Brand").
		Preload("ProductSkus", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
		Preload("ProductSkus.ComponentSpecs").
		Preload("ProductVariantOptions", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
		Preload("ProductVariantOptionValues", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
		Preload("ProductSkuVariants")

	if query != nil {
		queries = queries.Where("name ILIKE ?", fmt.Sprintf("%%%s%%", *query)) // ILIKE for case-insensitive in PostgreSQL
	}

	if productTypes != nil && len(*productTypes) > 0 {
		queries = queries.Where("product_type IN ?", *productTypes)
	}

	if productCategoryIds != nil && len(*productCategoryIds) > 0 {
		queries = queries.Where("product_category_id IN ?", *productCategoryIds)
	}

	if isActive != nil {
		queries = queries.Where("is_active = ?", *isActive)
	}

	if isShowOnlyInMarketPlace != nil {
		queries = queries.Where("is_show_only_in_market_place = ?", *isShowOnlyInMarketPlace)
	}

	// Hitung total
	queries.Count(&totalCount)

	if page != nil && pageSize != nil && *page > 0 && *pageSize > 0 {
		offset := (*page - 1) * *pageSize
		totalPage = (totalCount + int64(*pageSize) - 1) / int64(*pageSize)

		queries = queries.Limit(*pageSize).Offset(offset)
	}

	queries = queries.Order("products.modified_at DESC")
	queries.Find(&products)

	return products, totalCount, totalPage
}

func (repository *ProductRepositoryImpl) GetById(db *gorm.DB, productId uuid.UUID) (entity.Product, error) {
	var product entity.Product
	err := db.
		Preload("ProductCategory").
		Preload("Brand").
		Preload("ProductSkus", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
		Preload("ProductSkus.ComponentSpecs", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
		Preload("ProductSkus.ProductGroupItems").
		Preload("ProductSkus.ProductGroupItems.Product").
		Preload("ProductSkus.ProductGroupItems.Product.ProductSkus").
		Preload("ProductSkus.ProductGroupItems.Product.ProductSkus.ComponentSpecs").
		Preload("ProductSkus.ProductGroupItems.Product.ProductSkuVariants").
		Preload("ProductSkus.ProductGroupItems.Product.ProductVariantOptions").
		Preload("ProductSkus.ProductGroupItems.Product.ProductVariantOptionValues").
		Preload("ProductVariantOptions", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
		Preload("ProductVariantOptionValues", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
		Preload("ProductSkuVariants").
		First(&product, "id = ?", productId).Error
	return product, err
}

// Product SKU
func (repository *ProductRepositoryImpl) SaveProductSkus(db *gorm.DB, productSkus []entity.ProductSku) error {
	err := db.Save(&productSkus).Error
	return err
}

func (repository *ProductRepositoryImpl) UpdateStockProductSkus(db *gorm.DB, productSkus []entity.ProductSku) error {
	ids := []uuid.UUID{}
	caseStmt := "CASE id"
	for _, sku := range productSkus {
		ids = append(ids, sku.Id)
		if sku.Stock != nil {
			caseStmt += fmt.Sprintf(" WHEN '%s' THEN %d", sku.Id, *sku.Stock)
		} else {
			caseStmt += fmt.Sprintf(" WHEN '%s' THEN NULL", sku.Id)
		}
	}
	caseStmt += " END"

	query := fmt.Sprintf("UPDATE product_skus SET stock = %s WHERE id IN ?", caseStmt)

	err := db.Exec(query, ids).Error
	return err
}

func (repository *ProductRepositoryImpl) DeleteProductSkus(db *gorm.DB, productId uuid.UUID, productSkus []entity.ProductSku) error {
	err := db.Where("product_id = ?", productId).Delete(productSkus).Error
	return err
}

func (repository *ProductRepositoryImpl) BrowseProductSku(db *gorm.DB, query *string, productTypes *[]string, productCategoryId *uint, page, pageSize *int) ([]response.BrowseProductSku, int64, int64) {
	var productSkus []response.BrowseProductSku
	var totalCount int64 = 0
	var totalPage int64 = 0

	baseQuery := `
		FROM products p
		INNER JOIN product_skus ps ON p.id = ps.product_id
		WHERE p.is_active = true AND ps.is_active = true
	`
	var filters []string
	var params []interface{}

	if productCategoryId != nil {
		filters = append(filters, "p.product_category_id = ?")
		params = append(params, productCategoryId)
	}

	if query != nil && *query != "" {
		filters = append(filters, "p.name ILIKE ?")
		params = append(params, "%"+*query+"%")
	}

	if len(filters) > 0 {
		baseQuery += " AND " + strings.Join(filters, " AND ")
	}

	countQuery := "SELECT COUNT(*) " + baseQuery
	err := db.Raw(countQuery, params...).Scan(&totalCount).Error

	exception.PanicIfNeeded(err)

	dataQuery := `
		SELECT 
			ps.id,
			ps.product_id,
			ps.sku,
			p.name || ps.name AS name,
			ps.stock,
			ps.stock_alert,
			ps.price,
			p.is_active,
			--p.product_type,
			p.product_category_id,
			p.picture_id,
			p.description,
			p.created_at,
			p.modified_at
	` + baseQuery + `
		ORDER BY p.created_at, ps.sequence
	`

	if page == nil || pageSize == nil {
		err = db.Raw(dataQuery, params...).Scan(&productSkus).Error
		exception.PanicIfNeeded(err)

		return productSkus, totalCount, 1
	} else {

		dataQuery += " LIMIT ? OFFSET ?"

		offset := (*page - 1) * *pageSize
		paramsWithPagination := append(params, pageSize, offset)

		err = db.Raw(dataQuery, paramsWithPagination...).Scan(&productSkus).Error
		exception.PanicIfNeeded(err)

		totalPage = (totalCount + int64(*pageSize) - 1) / int64(*pageSize)

		return productSkus, totalCount, totalPage
	}
}

func (repository *ProductRepositoryImpl) GetProductSkuByIds(db *gorm.DB, productSkuIds []uuid.UUID) []response.BrowseProductSku {
	var productSkus []response.BrowseProductSku

	dataQuery := `
		SELECT 
			ps.id,
			ps.product_id,
			ps.sku,
			p.name || ps.name AS name,
			ps.stock,
			ps.stock_alert,
			ps.price,
			p.is_active,
			p.product_type,
			p.product_category_id,
			p.picture_id,
			p.description,
			p.created_at,
			p.modified_at
		FROM products p
		INNER JOIN product_skus ps 
			ON p.id = ps.product_id
		WHERE p.is_active = true AND ps.is_active = true
					AND ps.id in ?
		ORDER BY p.created_at, ps.sequence
	`
	db.Raw(dataQuery, productSkuIds).Scan(&productSkus)

	return productSkus
}

// func (repository *ProductRepositoryImpl) GetProductSkuByIds(db *gorm.DB, productSkuIds []uuid.UUID) []entity.ProductSku {
// 	var productSkus []entity.ProductSku

// 	db.Find(&productSkus, "id IN (?)", productSkuIds)
// 	return productSkus
// }

// Product SKU specs

func (repository *ProductRepositoryImpl) SaveComponentSpecs(db *gorm.DB, componentSpecs []entity.ComponentSpec) error {
	if len(componentSpecs) == 0 {
		return nil
	}
	err := db.Save(&componentSpecs).Error
	return err
}

func (repository *ProductRepositoryImpl) DeleteComponentSpecs(db *gorm.DB, productSkuId uuid.UUID, componentSpecs []entity.ComponentSpec) error {
	err := db.Where("ProductSkuId = ?", productSkuId).Delete(componentSpecs).Error
	return err
}

// Group Items
func (repository *ProductRepositoryImpl) SaveProductGroupItems(db *gorm.DB, groupItems []entity.ProductGroupItem) error {
	if len(groupItems) == 0 {
		return nil
	}
	err := db.Save(&groupItems).Error
	return err
}

func (repository *ProductRepositoryImpl) DeleteProductGroupItems(db *gorm.DB, parentId uuid.UUID, groupItems []entity.ProductGroupItem) error {
	err := db.Where("parent_id = ?", parentId).Delete(groupItems).Error
	return err
}

// Variant Options
func (repository *ProductRepositoryImpl) SaveProductVariantOptions(db *gorm.DB, variantOptions []entity.ProductVariantOption) error {
	if len(variantOptions) == 0 {
		return nil
	}

	err := db.Save(&variantOptions).Error
	return err
}

func (repository *ProductRepositoryImpl) DeleteProductVariantOptions(db *gorm.DB, productId uuid.UUID, variantOptions []entity.ProductVariantOption) error {
	err := db.Where("product_id = ?", productId).Delete(variantOptions).Error
	return err
}

// Variant Option Values
func (repository *ProductRepositoryImpl) SaveProductVariantOptionValues(db *gorm.DB, optionValues []entity.ProductVariantOptionValue) error {
	if len(optionValues) == 0 {
		return nil
	}
	err := db.Save(&optionValues).Error
	return err
}

func (repository *ProductRepositoryImpl) DeleteProductVariantOptionValues(db *gorm.DB, productId uuid.UUID, optionValues []entity.ProductVariantOptionValue) error {
	err := db.Where("product_id = ?", productId).Delete(optionValues).Error
	return err
}

// Product SKU Variant
func (repository *ProductRepositoryImpl) SaveProductSkuVariants(db *gorm.DB, productSkuVariants []entity.ProductSkuVariant) error {
	if len(productSkuVariants) == 0 {
		return nil
	}

	err := db.Save(&productSkuVariants).Error
	return err
}

func (repository *ProductRepositoryImpl) DeleteProductSkuVariants(db *gorm.DB, productId uuid.UUID, productSkuVariants []entity.ProductSkuVariant) error {
	err := db.Where("product_id = ?", productId).Delete(productSkuVariants).Error
	return err
}

// MD Variant Option.

func (repository *ProductRepositoryImpl) SaveVariantOption(db *gorm.DB, variantOption entity.VariantOption) (entity.VariantOption, error) {
	err := db.Save(&variantOption).Error
	return variantOption, err
}

func (repository *ProductRepositoryImpl) SaveVariantOptions(db *gorm.DB, variantOptions []entity.VariantOption) error {
	err := db.Save(&variantOptions).Error
	return err
}

func (repository *ProductRepositoryImpl) GetVariantOptions(db *gorm.DB) ([]entity.VariantOption, error) {
	var opt []entity.VariantOption
	err := db.Find(&opt).Error
	return opt, err
}
