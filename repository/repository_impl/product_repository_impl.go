package repository_impl

import (
	"fmt"
	"strings"
	"time"

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
	err := db.Omit(
		"ProductCategory",
		"Brand",
		"ProductSkus",
		// "ProductSkuDetails",
		"ProductVariantOptions",
		"ProductVariantOptionValues",
		"ProductSkuVariants",
		"ProductGroupItems",
	).Save(&product).Error
	return product, err
}

func (repository *ProductRepositoryImpl) Delete(db *gorm.DB, product entity.Product) error {
	err := db.Where("id = ?", product.Id).Update("deleted_at", time.Now().UTC()).Error
	return err
}

func (repository *ProductRepositoryImpl) Search(db *gorm.DB, query *string, productTypes *[]string, productCategoryId *uint, page, pageSize *int) ([]entity.Product, int64, int64) {
	var products []entity.Product
	var totalCount int64 = 0
	var totalPage int64 = 0
	var offset int = 0

	queries := db.Model(&products).
		Preload("ProductSkus", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
		Preload("ProductSkus.ProductSpecs").
		Preload("ProductVariantOptions", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
		Preload("ProductVariantOptionValues", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
		Preload("ProductSkuVariants")

	if query != nil {
		queries.Where("name like ?", "%"+*query+"%")
	}

	if productTypes != nil {
		placeholders := make([]string, len(*productTypes))
		for i := range *productTypes {
			placeholders[i] = "?"
		}
		inClause := strings.Join(placeholders, ",")
		queries.Where(fmt.Sprintf("product_type in (%s)", inClause), repository.interfaceSlice(*productTypes)...)
	}

	if productCategoryId != nil {
		queries.Where("product_category_id = ?", productCategoryId)
	}

	queries.Count(&totalCount)
	if page != nil && pageSize != nil {
		offset = (*page - 1) * *pageSize
		totalPage = ((totalCount + int64(*pageSize) - 1) / int64(*pageSize))

		queries.Limit(*pageSize).Offset(offset).Order("products.modified_at desc").Find(&products)
	} else {
		queries.Order("products.modified_at desc").Find(&products)
	}
	return products, totalCount, totalPage
}

func (repository *ProductRepositoryImpl) BrowseProductSku(db *gorm.DB, query *string, productTypes *[]string, productCategoryId *uint, page, pageSize *int) ([]response.BrowseProductSku, int64, int64) {
	var products []response.BrowseProductSku
	var totalCount int64 = 0
	var totalPage int64 = 0
	// var offset int = 0

	// queries := db.Model(&products).
	// 	Preload("ProductSkus", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
	// 	// Preload("ProductSkus.ProductSkuDetails").
	// 	Preload("ProductVariantOptions", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
	// 	Preload("ProductVariantOptionValues", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
	// 	Preload("ProductSkuVariants")

	db.Raw(`
		Select 
			ps.id,
			ps.product_id,
			ps.sku,
			p.name || ps.name as name,
			ps.stock,
			ps.stock_alert,
			ps.price,
			p.is_active,
			p.product_type,
			p.picture_id,
			p.description,
			p.created_at,
			p.modified_at
		from products p
		INNER JOIN product_skus ps
				ON p.id = ps.product_id

		ORDER BY p.created_at, sequence
	`).Find(&products)

	return products, totalCount, totalPage
}

// func (repository *ProductRepositoryImpl) SearchProductSku(
// 	db *gorm.DB,
// 	companyId uuid.UUID,
// 	outletId *uuid.UUID,
// 	query *string,
// 	productTypes *[]string,
// 	isInventoryOnly *bool,
// 	productCategoryId *uuid.UUID,
// 	brandId *uuid.UUID,
// 	page, pageSize *int,
// ) ([]db_view.ProductSkuView, int64, int64) {
// 	var result []db_view.ProductSkuView

// 	rows, err := repository.sqlDb.Query(
// 		"EXEC ProductSku_Get @CompanyId = @CompanyId, @OutletId = @OutletId",
// 		sql.Named("CompanyId", companyId),
// 		sql.Named("OutletId", outletId),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer rows.Close()

// 	// Map product skus
// 	productSkuMap := make(map[string]*db_view.ProductSkuView)

// 	for rows.Next() {
// 		var sku db_view.ProductSkuView

// 		err := rows.Scan(
// 			&sku.Id,
// 			&sku.ProductId,
// 			&sku.ProductType,
// 			&sku.SKU,
// 			&sku.ParentSKU,
// 			&sku.Name,
// 			&sku.IsPartOfCompositeOnly,
// 			&sku.IsTrackInventory,
// 			&sku.PictureId,
// 			&sku.ModifiedAt,
// 		)
// 		if err != nil {
// 			panic(err)
// 		}

// 		// Use string representation of GUID as map key
// 		key := sku.Id.String()
// 		productSkuMap[key] = &sku
// 	}

// 	// Move to next result set (ProductSkuDetails)
// 	if !rows.NextResultSet() {
// 		return result, 0, 0 // No second result set
// 	}

// 	for rows.Next() {
// 		var detail entity.ProductSkuDetail
// 		var productSkuId uuid.UUID

// 		err := rows.Scan(
// 			&detail.Id,
// 			&productSkuId,
// 			&detail.OutletId,
// 			&detail.Fee,
// 			&detail.CostPrice,
// 			&detail.SellingPrice,
// 			&detail.Qty,
// 			&detail.QtyAlert,
// 		)
// 		if err != nil {
// 			panic(err)
// 		}

// 		key := productSkuId.String()
// 		if sku, ok := productSkuMap[key]; ok {
// 			sku.ProductSkuDetails = append(sku.ProductSkuDetails, detail)
// 		}
// 	}

// 	// Flatten the map into a slice
// 	for _, sku := range productSkuMap {
// 		result = append(result, *sku)
// 	}

// 	totalData := int64(len(result))
// 	totalPages := int64(1)
// 	if pageSize != nil && *pageSize > 0 {
// 		totalPages = (totalData + int64(*pageSize) - 1) / int64(*pageSize)
// 	}

// 	return result, totalData, totalPages
// }

func (repository *ProductRepositoryImpl) GetById(db *gorm.DB, productId uuid.UUID) (entity.Product, error) {
	var product entity.Product
	err := db.
		Preload("ProductCategory").
		Preload("ProductSkus", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
		Preload("ProductSkus.ProductSPecs", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
		Preload("ProductGroupItems").
		Preload("ProductGroupItems.Product").
		Preload("ProductGroupItems.Product.ProductSkus").
		Preload("ProductGroupItems.Product.ProductSkus.ProductSpecs").
		Preload("ProductGroupItems.Product.ProductSkuVariants").
		Preload("ProductGroupItems.Product.ProductVariantOptions").
		Preload("ProductGroupItems.Product.ProductVariantOptionValues").
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

func (repository *ProductRepositoryImpl) DeleteProductSkus(db *gorm.DB, productId uuid.UUID, productSkus []entity.ProductSku) error {
	err := db.Where("product_id = ?", productId).Delete(productSkus).Error
	return err
}

// Product SKU Details

func (repository *ProductRepositoryImpl) SaveProductSpecs(db *gorm.DB, productSpecs []entity.ProductSpec) error {
	err := db.Save(&productSpecs).Error
	return err
}

func (repository *ProductRepositoryImpl) DeleteProductSpecs(db *gorm.DB, productSkuId uuid.UUID, productSpecs []entity.ProductSpec) error {
	err := db.Where("ProductSkuId = ?", productSkuId).Delete(productSpecs).Error
	return err
}

// Group Items
func (repository *ProductRepositoryImpl) SaveProductGroupItems(db *gorm.DB, groupItems []entity.ProductGroupItem) error {
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

// Helpers
func (repository *ProductRepositoryImpl) interfaceSlice(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}
