package service_impl

import (
	"fmt"
	"slices"
	"time"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers"
	"api.mijkomp.com/models/entity"
	"api.mijkomp.com/models/enum"
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductServiceImpl struct {
	ProductRepository repository.ProductRepository

	db *gorm.DB
}

func NewProductService(userRepostitory repository.ProductRepository, db *gorm.DB) *ProductServiceImpl {
	return &ProductServiceImpl{
		ProductRepository: userRepostitory,

		db: db,
	}
}

func (service *ProductServiceImpl) Create(currentUserId uint, productId uuid.UUID, payload request.ProductPayload) response.ProductResponse {

	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	product := entity.Product{
		Id:                productId,
		ProductType:       payload.ProductType,
		SKU:               payload.SKU,
		Name:              payload.Name,
		IsActive:          payload.IsActive,
		PictureId:         payload.PictureId,
		ProductCategoryId: payload.ProductCategoryId,
		Description:       payload.Description,
		CreatedById:       currentUserId,
		CreatedAt:         time.Now().UTC(),
		ModifiedById:      currentUserId,
		ModifiedAt:        time.Now().UTC(),
	}

	productRes, err := service.ProductRepository.Save(tx, product)
	exception.PanicIfNeeded(err)

	// Product SKU
	productSkus := []entity.ProductSku{}
	// productSkuDetails := []entity.ProductSkuDetail{}

	for i, productSku := range payload.ProductSkus {
		productSkus = append(productSkus, entity.ProductSku{
			Id:        productSku.Id,
			ProductId: productId,
			SKU:       productSku.SKU,
			Sequence:  i + 1,
			IsActive:  true,
		})

		// for _, skuDetail := range productSku.ProductSkuDetails {
		// 	productSkuDetail := entity.ProductSkuDetail{
		// 		Id:           skuDetail.Id,
		// 		ProductSkuId: productSku.Id,
		// 		OutletId:     skuDetail.OutletId,
		// 		Fee:          skuDetail.Fee,
		// 		CostPrice:    skuDetail.CostPrice,
		// 		SellingPrice: skuDetail.SellingPrice,
		// 		Qty:          skuDetail.Qty,
		// 		QtyAlert:     skuDetail.QtyAlert,
		// 		CreatedById:  currentUserId,
		// 		CreatedAt:    time.Now().UTC(),
		// 		ModifiedById: currentUserId,
		// 		ModifiedAt:   time.Now().UTC(),
		// 	}

		// 	if !product.IsTrackInventory {
		// 		productSkuDetail.Qty = nil
		// 		productSkuDetail.QtyAlert = nil
		// 	}

		// 	productSkuDetails = append(productSkuDetails, productSkuDetail)
		// }
	}

	err = service.ProductRepository.SaveProductSkus(tx, productSkus)
	exception.PanicIfNeeded(err)

	// err = service.ProductRepository.SaveProductSkuDetails(tx, productSkuDetails)
	// exception.PanicIfNeeded(err)

	// Product Variant
	if payload.ProductType == enum.ProductTypeVariant {

		productVariantOptions := []entity.ProductVariantOption{}
		productVariantOptionValues := []entity.ProductVariantOptionValue{}
		productSkuVariants := []entity.ProductSkuVariant{}

		for i, variantOption := range payload.ProductVariantOptions {
			productVariantOptions = append(productVariantOptions, entity.ProductVariantOption{
				Id:        variantOption.Id,
				ProductId: productId,
				Name:      variantOption.Name,
				Sequence:  i + 1,
			})

		}

		for i, variantOptionValue := range payload.ProductVariantOptionValues {
			productVariantOptionValues = append(productVariantOptionValues, entity.ProductVariantOptionValue{
				Id:                     variantOptionValue.Id,
				ProductId:              productId,
				ProductVariantOptionId: variantOptionValue.ProductVariantOptionId,
				Name:                   variantOptionValue.Name,
				Sequence:               i + 1,
			})
		}

		for _, productSkuVariant := range payload.ProductSkuVariants {
			productVariantOptionIdx := slices.IndexFunc(productVariantOptions, func(opt entity.ProductVariantOption) bool { return opt.Id == productSkuVariant.ProductVariantOptionId })

			if productVariantOptionIdx != -1 {
				productSkuVariants = append(productSkuVariants, entity.ProductSkuVariant{
					Id:                          productSkuVariant.Id,
					ProductId:                   productId,
					ProductSkuId:                productSkuVariant.ProductSkuId,
					ProductVariantOptionId:      productSkuVariant.ProductVariantOptionId,
					ProductVariantOptionValueId: productSkuVariant.ProductVariantOptionValueId,
					Sequence:                    productVariantOptions[productVariantOptionIdx].Sequence,
				})
			}
		}

		err = service.ProductRepository.SaveProductVariantOptions(tx, productVariantOptions)
		exception.PanicIfNeeded(err)

		err = service.ProductRepository.SaveProductVariantOptionValues(tx, productVariantOptionValues)
		exception.PanicIfNeeded(err)

		err = service.ProductRepository.SaveProductSkus(tx, productSkus)
		exception.PanicIfNeeded(err)

		err = service.ProductRepository.SaveProductSkuVariants(tx, productSkuVariants)
		exception.PanicIfNeeded(err)
	}

	// Product Group
	if payload.ProductType == enum.ProductTypeGroup {
		productGroupItems := []entity.ProductGroupItem{}
		for _, el := range payload.ProductGroupItems {
			productGroupItems = append(productGroupItems, entity.ProductGroupItem{
				Id:           el.Id,
				ParentId:     productRes.Id,
				ProductId:    el.ProductId,
				ProductSkuId: el.ProductSkuId,
				Qty:          el.Qty,
			})
		}

		if len(productGroupItems) > 0 {
			err := service.ProductRepository.SaveProductGroupItems(tx, productGroupItems)
			exception.PanicIfNeeded(err)
		} else {
			panic(exception.NewValidationError("Group Detail tidak boleh kosong"))
		}
	}

	res, err := service.ProductRepository.GetById(tx, productRes.Id)
	exception.PanicIfNeeded(err)

	return service.mapProduct(res)
}

func (service *ProductServiceImpl) Update(currentUserId uint, productId uuid.UUID, payload request.ProductPayload) response.ProductResponse {
	// tx
	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	product, err := service.ProductRepository.GetById(service.db, productId)
	exception.PanicIfNeeded(err)

	product.ProductType = payload.ProductType
	product.Name = payload.Name
	product.SKU = payload.SKU
	product.IsActive = payload.IsActive
	product.PictureId = payload.PictureId
	product.ProductCategoryId = payload.ProductCategoryId
	product.Description = payload.Description
	product.ModifiedById = currentUserId
	product.ModifiedAt = time.Now().UTC()

	_, err = service.ProductRepository.Save(tx, product)
	exception.PanicIfNeeded(err)

	// Product SKU (Update, Delete, Add New)
	productSkus := product.ProductSkus
	productSkuToBeDeleted := []entity.ProductSku{}

	// productSkuDetails := []entity.ProductSkuDetail{}
	// productSkuDetailsToBeDeleted := []entity.ProductSkuDetail{}
	// Update, delete
	for i, productSku := range productSkus {
		payloadIdx := slices.IndexFunc(payload.ProductSkus, func(model request.ProductSkuPayload) bool {
			return model.Id == productSku.Id
		})

		if payloadIdx != -1 {
			productSkus[i].SKU = payload.ProductSkus[payloadIdx].SKU
			productSkus[i].Name = payload.ProductSkus[payloadIdx].Name
			productSkus[i].IsActive = payload.ProductSkus[payloadIdx].IsActive
		} else {
			productSkuToBeDeleted = append(productSkuToBeDeleted, productSku)
		}
	}

	// Add new
	for _, productSku := range payload.ProductSkus {
		savedIdx := slices.IndexFunc(productSkus, func(model entity.ProductSku) bool {
			return model.Id == productSku.Id
		})

		if savedIdx == -1 {
			productSkus = append(productSkus, entity.ProductSku{
				Id:        productSku.Id,
				ProductId: product.Id,
				SKU:       productSku.SKU,
				Sequence:  len(productSkus) + 1,
			})
		}
	}

	err = service.ProductRepository.DeleteProductSkus(tx, product.Id, productSkuToBeDeleted)
	exception.PanicIfNeeded(err)

	err = service.ProductRepository.SaveProductSkus(tx, productSkus)
	exception.PanicIfNeeded(err)

	// // Product Sku Detail
	// err = service.ProductRepository.DeleteProductSkuDetails(tx, product.Id, productSkuDetailsToBeDeleted)
	// exception.PanicIfNeeded(err)

	// err = service.ProductRepository.SaveProductSkuDetails(tx, productSkuDetails)
	// exception.PanicIfNeeded(err)

	if product.ProductType == enum.ProductTypeVariant {
		// Product variant options (Update, Delete, Add New)
		productVariantOptions := product.ProductVariantOptions
		productVariantOptionsToBeDeleted := []entity.ProductVariantOption{}

		for i, variantOption := range productVariantOptions {
			optionIdx := slices.IndexFunc(payload.ProductVariantOptions, func(model request.ProductVariantOptionPayload) bool {
				return model.Id == variantOption.Id
			})

			if optionIdx != -1 {

				productVariantOptions[i].Name = payload.ProductVariantOptions[optionIdx].Name
				productVariantOptions[i].Sequence = variantOption.Sequence
			} else {
				productVariantOptionsToBeDeleted = append(productVariantOptionsToBeDeleted, variantOption)
			}
		}

		// Add new
		for _, variantOption := range payload.ProductVariantOptions {
			savedIdx := slices.IndexFunc(productVariantOptions, func(model entity.ProductVariantOption) bool {
				return model.Id == variantOption.Id
			})

			if savedIdx == -1 {
				productVariantOptions = append(productVariantOptions, entity.ProductVariantOption{
					Id:        variantOption.Id,
					ProductId: product.Id,
					Name:      variantOption.Name,
					Sequence:  len(productVariantOptions) + 1,
				})
			}
		}

		err = service.ProductRepository.DeleteProductVariantOptions(tx, product.Id, productVariantOptionsToBeDeleted)
		exception.PanicIfNeeded(err)

		err = service.ProductRepository.SaveProductVariantOptions(tx, productVariantOptions)
		exception.PanicIfNeeded(err)

		// Product variant option values (Update, Delete, Add New)
		productVariantOptionValues := product.ProductVariantOptionValues
		productVariantOptionValuesToBeDeleted := []entity.ProductVariantOptionValue{}

		for i, variantOptionValue := range productVariantOptionValues {
			optionValueIdx := slices.IndexFunc(payload.ProductVariantOptionValues, func(model request.ProductVariantOptionValuePayload) bool {
				return model.Id == variantOptionValue.Id
			})

			if optionValueIdx != -1 {
				productVariantOptionValues[i].ProductVariantOptionId = payload.ProductVariantOptionValues[optionValueIdx].ProductVariantOptionId
				productVariantOptionValues[i].Name = payload.ProductVariantOptionValues[optionValueIdx].Name
				productVariantOptionValues[i].Sequence = variantOptionValue.Sequence
			} else {
				productVariantOptionValuesToBeDeleted = append(productVariantOptionValuesToBeDeleted, variantOptionValue)
			}
		}

		// Add new
		for _, variantOptionValue := range payload.ProductVariantOptionValues {
			savedIdx := slices.IndexFunc(productVariantOptionValues, func(model entity.ProductVariantOptionValue) bool {
				return model.Id == variantOptionValue.Id
			})

			if savedIdx == -1 {
				productVariantOptionValues = append(productVariantOptionValues, entity.ProductVariantOptionValue{
					Id:                     variantOptionValue.Id,
					ProductId:              product.Id,
					ProductVariantOptionId: variantOptionValue.ProductVariantOptionId,
					Name:                   variantOptionValue.Name,
					Sequence:               len(productVariantOptionValues) + 1,
				})
			}
		}

		err = service.ProductRepository.DeleteProductVariantOptionValues(tx, productId, productVariantOptionValuesToBeDeleted)
		exception.PanicIfNeeded(err)

		err = service.ProductRepository.SaveProductVariantOptionValues(tx, productVariantOptionValues)
		exception.PanicIfNeeded(err)

		// Product SKU variant (Update, Delete, Add New)
		productSkuVariants := product.ProductSkuVariants
		productSkuVariantsToBeDeleted := []entity.ProductSkuVariant{}

		for i, skuVariant := range productSkuVariants {
			skuVariantIdx := slices.IndexFunc(payload.ProductSkuVariants, func(model request.ProductSkuVariantPayload) bool {
				return model.Id == skuVariant.Id
			})

			if skuVariantIdx != -1 {
				productSkuVariants[i].ProductSkuId = payload.ProductSkuVariants[skuVariantIdx].ProductSkuId
				productSkuVariants[i].ProductVariantOptionId = payload.ProductSkuVariants[skuVariantIdx].ProductVariantOptionId
				productSkuVariants[i].ProductVariantOptionValueId = payload.ProductSkuVariants[skuVariantIdx].ProductVariantOptionValueId
			} else {
				productSkuVariantsToBeDeleted = append(productSkuVariantsToBeDeleted, skuVariant)
			}
		}

		// Add new
		for _, skuVariant := range payload.ProductSkuVariants {
			savedIdx := slices.IndexFunc(productSkuVariants, func(model entity.ProductSkuVariant) bool {
				return model.Id == skuVariant.Id
			})

			if savedIdx == -1 {
				productSkuVariants = append(productSkuVariants, entity.ProductSkuVariant{
					Id:                          skuVariant.Id,
					ProductId:                   product.Id,
					ProductSkuId:                skuVariant.ProductSkuId,
					ProductVariantOptionId:      skuVariant.ProductVariantOptionId,
					ProductVariantOptionValueId: skuVariant.ProductVariantOptionValueId,
					Sequence:                    len(productVariantOptionValues) + 1,
				})
			}
		}

		err = service.ProductRepository.DeleteProductSkuVariants(tx, productId, productSkuVariantsToBeDeleted)
		exception.PanicIfNeeded(err)

		err = service.ProductRepository.SaveProductSkuVariants(tx, productSkuVariants)
		exception.PanicIfNeeded(err)
	}

	// Product Group Items (Update, Delete, Add New)
	if product.ProductType == enum.ProductTypeGroup {
		if len(payload.ProductGroupItems) == 0 {
			panic(exception.NewValidationError("Group detail tidak boleh kosong"))
		}

		productGroupItems := product.ProductGroupItems
		productGroupItemsToBeDeleted := []entity.ProductGroupItem{}

		for i, groupItem := range productGroupItems {
			groupItemIdx := slices.IndexFunc(payload.ProductGroupItems, func(model request.ProductGroupItemPayload) bool {
				return model.Id == groupItem.Id
			})

			if groupItemIdx != -1 {
				productGroupItems[i].ParentId = groupItem.ParentId
				productGroupItems[i].ProductId = groupItem.ProductId
				productGroupItems[i].ProductSkuId = groupItem.ProductSkuId
				productGroupItems[i].Qty = payload.ProductGroupItems[groupItemIdx].Qty
			} else {
				productGroupItemsToBeDeleted = append(productGroupItemsToBeDeleted, groupItem)
				toBeDeletedIdx := slices.IndexFunc(productGroupItems, func(model entity.ProductGroupItem) bool {
					return model.Id == groupItem.Id
				})
				productGroupItems = append(productGroupItems[:toBeDeletedIdx], productGroupItems[toBeDeletedIdx+1:]...)
			}
		}

		// Add new
		for _, groupItem := range payload.ProductGroupItems {
			savedIdx := slices.IndexFunc(productGroupItems, func(model entity.ProductGroupItem) bool {
				return model.Id == groupItem.Id
			})

			if savedIdx == -1 {
				productGroupItems = append(productGroupItems, entity.ProductGroupItem{
					Id:           groupItem.Id,
					ParentId:     product.Id,
					ProductId:    groupItem.ProductId,
					ProductSkuId: groupItem.ProductSkuId,
					Qty:          groupItem.Qty,
				})
			}
		}

		err = service.ProductRepository.DeleteProductGroupItems(tx, productId, productGroupItemsToBeDeleted)
		exception.PanicIfNeeded(err)

		err = service.ProductRepository.SaveProductGroupItems(tx, productGroupItems)
		exception.PanicIfNeeded(err)
	}
	// return new result
	res, err := service.ProductRepository.GetById(tx, productId)
	exception.PanicIfNeeded(err)

	return service.mapProduct(res)
}

func (service *ProductServiceImpl) Delete(currentUserId uint, productId uuid.UUID) string {
	// tx
	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	product, err := service.ProductRepository.GetById(service.db, productId)
	exception.PanicIfNeeded(err)

	err = service.ProductRepository.Delete(service.db, product)
	exception.PanicIfNeeded(err)

	return fmt.Sprintf("Produk %s berhasil di hapus", product.Name)
}

func (service *ProductServiceImpl) Search(currentUserId uint, query *string, productTypes *[]string, productCategoryId *uint, page, pageSize *int) response.PageResult {

	res, totalCount, totalPage := service.ProductRepository.Search(service.db, query, productTypes, productCategoryId, page, pageSize)

	return response.PageResult{
		Items:      service.mapProducts(res),
		TotalCount: totalCount,
		PageSize:   totalPage,
	}
}

// func (service *ProductServiceImpl) SearchProductSku(currentUserId uint, outletId *uuid.UUID, query *string, productTypes *[]string, isInventoryOnly *bool, productCategoryId *uuid.UUID, brandId *uuid.UUID, page, pageSize *int) response.PageResult {

// 	res, totalCount, totalPage := service.ProductRepository.SearchProductSku(service.db, outletId, query, productTypes, isInventoryOnly, productCategoryId, brandId, page, pageSize)

// 	return response.PageResult{
// 		Items:      service.mapBrowseProductSku(res),
// 		TotalCount: totalCount,
// 		PageSize:   totalPage,
// 	}
// }

func (service *ProductServiceImpl) GetById(currentUserId uint, productId uuid.UUID) response.ProductResponse {
	res, err := service.ProductRepository.GetById(service.db, productId)
	exception.PanicIfNeeded(err)

	return service.mapProduct(res)
}

// Variant Options
func (service *ProductServiceImpl) CreateVariantOptions(currentUserId uint, payload request.VariantOptionPayload) response.VariantOption {

	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	variantOption := entity.VariantOption{
		Name: payload.Name,
	}

	res, err := service.ProductRepository.SaveVariantOption(tx, variantOption)
	exception.PanicIfNeeded(err)

	return service.mapVariantOption(res)
}

func (service *ProductServiceImpl) GetVariantOptions(currentUserId uint) []response.VariantOption {

	res, err := service.ProductRepository.GetVariantOptions(service.db)
	exception.PanicIfNeeded(err)

	return service.mapVariantOptions(res)
}

// Map helpers
func (service *ProductServiceImpl) mapProducts(products []entity.Product) []response.ProductResponse {
	productRes := []response.ProductResponse{}

	for _, el := range products {
		productRes = append(productRes, service.mapProduct(el))
	}

	return productRes
}

func (service *ProductServiceImpl) mapProduct(product entity.Product) response.ProductResponse {
	productRes := response.ProductResponse{
		Id:           product.Id,
		SKU:          product.SKU,
		Name:         product.Name,
		IsActive:     product.IsActive,
		PictureId:    product.PictureId,
		Description:  product.Description,
		CreatedById:  product.CreatedById,
		CreatedAt:    product.CreatedAt,
		ModifiedById: product.ModifiedById,
		ModifiedAt:   product.ModifiedAt,

		ProductSkus:                []response.ProductSku{},
		ProductGroupItems:          []response.ProductGroupItemResponse{},
		ProductVariantOptions:      []response.ProductVariantOption{},
		ProductVariantOptionValues: []response.ProductVariantOptionValue{},
		ProductSkuVariants:         []response.ProductSkuVariant{},
	}
	productRes.ProductType = response.EnumResponse{
		Code: string(product.ProductType),
		Name: product.ProductType.DisplayString(),
	}

	if product.ProductCategory != nil {
		productCategoryRes := response.ProductCategoryResponse{
			Id:   product.ProductCategory.Id,
			Name: product.ProductCategory.Name,
		}
		productRes.ProductCategory = &productCategoryRes
	}

	for _, productSku := range product.ProductSkus {
		productSkuRes := response.ProductSku{
			Id:        productSku.Id,
			ProductId: productSku.ProductId,
			SKU:       productSku.SKU,
			Name:      productSku.Name,
			IsActive:  productSku.IsActive,
			Sequence:  productSku.Sequence,
			// ProductSkuDetails: service.mapProductSkuDetails(productSku.ProductSkuDetails),
		}

		productRes.ProductSkus = append(productRes.ProductSkus, productSkuRes)
	}

	if product.ProductType == enum.ProductTypeVariant {

		for _, variantOption := range product.ProductVariantOptions {
			productRes.ProductVariantOptions = append(productRes.ProductVariantOptions, response.ProductVariantOption{
				Id:          variantOption.Id,
				ProductId:   variantOption.ProductId,
				Name:        variantOption.Name,
				AllowDelete: false,
				Sequence:    variantOption.Sequence,
			})
		}

		for _, variantOptionValue := range product.ProductVariantOptionValues {
			productRes.ProductVariantOptionValues = append(productRes.ProductVariantOptionValues, response.ProductVariantOptionValue{
				Id:                     variantOptionValue.Id,
				ProductVariantOptionId: variantOptionValue.ProductVariantOptionId,
				Name:                   variantOptionValue.Name,
				AllowDelete:            false,
				Sequence:               variantOptionValue.Sequence,
			})
		}

		for _, productSkuVariant := range product.ProductSkuVariants {
			productRes.ProductSkuVariants = append(productRes.ProductSkuVariants, response.ProductSkuVariant{
				Id:                          productSkuVariant.Id,
				ProductSkuId:                productSkuVariant.ProductSkuId,
				ProductVariantOptionId:      productSkuVariant.ProductVariantOptionId,
				ProductVariantOptionValueId: productSkuVariant.ProductVariantOptionValueId,
			})
		}

	} else if product.ProductType == enum.ProductTypeGroup {

		for _, item := range product.ProductGroupItems {

			productRes.ProductGroupItems = append(productRes.ProductGroupItems, response.ProductGroupItemResponse{
				Id:           item.Id,
				ProductId:    item.Product.Id,
				ProductSkuId: item.ProductSkuId,
				Qty:          item.Qty,
				Product:      service.mapProduct(item.Product),
			})

			// if item.Product.ProductType == enum.ProductTypeSimple {

			// 	productRes.ProductGroupItems = append(productRes.ProductGroupItems, response.ProductGroupItemResponse{
			// 		Id:           item.Id,
			// 		ProductId:    item.Product.Id,
			// 		ProductSkuId: item.ProductSkuId,
			// 		Qty:          item.Qty,
			// 		Product:      service.mapProduct(item.Product),
			// 	})

			// } else if item.Product.ProductType == enum.ProductTypeVariant {
			// 	productSku := entity.ProductSku{}
			// 	skuVariants := []entity.ProductSkuVariant{}
			// 	optionValues := []entity.ProductVariantOptionValue{}

			// 	for _, sku := range item.Product.ProductSkus {
			// 		if sku.Id == item.ProductSkuId {
			// 			productSku = sku
			// 		}
			// 	}

			// 	for _, skuVariant := range item.Product.ProductSkuVariants {
			// 		if skuVariant.ProductId == item.ProductId {
			// 			skuVariants = append(skuVariants, skuVariant)

			// 			// for _, optValue := range item.Product.ProductVariantOptionValues {
			// 			// 	if skuVariant.ProductVariantOptionValueId == optValue.Id {
			// 			// 		optionValues = append(optionValues, optValue)
			// 			// 	}
			// 			// }
			// 		}
			// 	}

			// 	newGroupItem := response.ProductGroupItemResponse{
			// 		Id:           item.Id,
			// 		ProductId:    item.Product.Id,
			// 		ProductSkuId: item.ProductSkuId,
			// 		ProductType:  item.Product.ProductType,
			// 		ProductSKU:   productSku.SKU,
			// 		ProductName:  item.Product.Name,
			// 		// ProductCostPrice:    productSku.CostPrice,
			// 		// ProductSellingPrice: productSku.SellingPrice,
			// 		Qty: item.Qty,

			// 		Product: service.mapProduct(item.Product),
			// 	}

			// 	for _, skuVariant := range skuVariants {
			// 		optionValueIdx := slices.IndexFunc(optionValues, func(model entity.ProductVariantOptionValue) bool {
			// 			return model.Id == skuVariant.ProductVariantOptionValueId
			// 		})

			// 		if optionValueIdx != -1 {
			// 			newGroupItem.ProductName += fmt.Sprintf(" %s / %s ", newGroupItem.ProductName, optionValues[optionValueIdx].Name)
			// 		}
			// 	}

			// 	productRes.ProductGroupItems = append(productRes.ProductGroupItems, newGroupItem)

			// }

		}

	}

	return productRes
}

// Variant Options Map Helpers

// func (service *ProductServiceImpl) mapProductSkuDetails(skuDetails []entity.ProductSkuDetail) []response.ProductSkuDetail {

// 	productSkuDetailRes := []response.ProductSkuDetail{}
// 	for _, skuDetail := range skuDetails {
// 		productSkuDetailRes = append(productSkuDetailRes, response.ProductSkuDetail{
// 			Id:           skuDetail.Id,
// 			ProductSkuId: skuDetail.ProductSkuId,
// 			OutletId:     skuDetail.OutletId,
// 			Fee:          skuDetail.Fee,
// 			CostPrice:    skuDetail.CostPrice,
// 			SellingPrice: skuDetail.SellingPrice,
// 			Qty:          skuDetail.Qty,
// 			QtyAlert:     skuDetail.QtyAlert,
// 		})
// 	}
// 	return productSkuDetailRes
// }

func (service *ProductServiceImpl) mapVariantOptions(opts []entity.VariantOption) []response.VariantOption {
	result := []response.VariantOption{}

	for _, item := range opts {
		result = append(result, service.mapVariantOption(item))
	}

	return result
}

func (service *ProductServiceImpl) mapVariantOption(opt entity.VariantOption) response.VariantOption {
	return response.VariantOption{
		Id:   opt.Id,
		Name: opt.Name,
	}
}

// func (service *ProductServiceImpl) mapBrowseProductSku(productSkus []db_view.ProductSkuView) []response.BrowseProductSkuResponse {
// 	res := []response.BrowseProductSkuResponse{}
// 	for _, item := range productSkus {
// 		res = append(res, response.BrowseProductSkuResponse{
// 			Id:                      item.Id,
// 			ProductId:               item.ProductId,
// 			Name:                    item.Name,
// 			SKU:                     item.SKU,
// 			ParentSKU:               item.ParentSKU,
// 			IsActive:                item.IsActive,
// 			IsPartOfCompositeOnly:   item.IsPartOfCompositeOnly,
// 			IsTrackInventory:        item.IsTrackInventory,
// 			IsPriceSameForAllOutlet: item.IsPriceSameForAllOutlet,
// 			PictureId:               item.PictureId,
// 			Description:             item.Description,
// 			CreatedById:             item.CreatedById,
// 			CreatedAt:               item.CreatedAt,
// 			ModifiedById:            item.ModifiedById,
// 			ModifiedAt:              item.ModifiedAt,
// 			ProductSkuDetails:       service.mapProductSkuDetails(item.ProductSkuDetails),
// 		})
// 	}

// 	return res
// }
