package service_impl

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers"
	"api.mijkomp.com/helpers/data_mapper"
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
	db                *gorm.DB
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
		Id:                      productId,
		ProductType:             payload.ProductType,
		SKU:                     payload.SKU,
		Name:                    payload.Name,
		IsActive:                true,
		IsShowOnlyInMarketPlace: payload.IsShowOnlyInMarketPlace,
		ImageIds:                helpers.JoinImageIds(payload.ImageIds),
		VideoUrl:                payload.VideoUrl,
		Tags:                    payload.Tags,
		ProductCategoryId:       payload.ProductCategoryId,
		BrandId:                 payload.BrandId,
		ComponentTypeId:         payload.ComponentTypeId,
		Description:             payload.Description,
		ProductSpec:             payload.ProductSpec,
		CreatedById:             currentUserId,
		CreatedAt:               time.Now().UTC(),
		ModifiedById:            currentUserId,
		ModifiedAt:              time.Now().UTC(),
	}

	productRes, err := service.ProductRepository.Save(tx, product)
	exception.PanicIfNeeded(err)

	// Product SKU
	productSkus := []entity.ProductSku{}
	componentSpecs := []entity.ComponentSpec{}
	productGroupItems := []entity.ProductGroupItem{}

	for i, productSku := range payload.ProductSkus {
		productSkus = append(productSkus, entity.ProductSku{
			Id:         productSku.Id,
			ProductId:  productId,
			Name:       productSku.Name,
			SKU:        productSku.SKU,
			Price:      productSku.Price,
			Stock:      productSku.Stock,
			StockAlert: productSku.StockAlert,
			Sequence:   i + 1,
			IsActive:   true,
		})

		for i, componentSpec := range productSku.ComponentSpecs {
			componentSpec := entity.ComponentSpec{
				Id:           componentSpec.Id,
				ProductSkuId: productSku.Id,
				SpecKey:      componentSpec.SpecKey,
				SpecValue:    componentSpec.SpecValue,
				Sequence:     i + 1,
			}

			componentSpecs = append(componentSpecs, componentSpec)
		}

		// Product Group
		if payload.ProductType == enum.ProductTypeGroup {
			for _, el := range productSku.ProductGroupItems {
				productGroupItems = append(productGroupItems, entity.ProductGroupItem{
					Id:           el.Id,
					ParentId:     productSku.Id,
					ProductId:    el.ProductId,
					ProductSkuId: el.ProductSkuId,
					Qty:          el.Qty,
				})
			}
		}
	}

	err = service.ProductRepository.SaveProductSkus(tx, productSkus)
	exception.PanicIfNeeded(err)

	err = service.ProductRepository.SaveComponentSpecs(tx, componentSpecs)
	exception.PanicIfNeeded(err)

	if payload.ProductType == enum.ProductTypeGroup {
		if len(productGroupItems) > 0 {
			err = service.ProductRepository.SaveProductGroupItems(tx, productGroupItems)
			exception.PanicIfNeeded(err)
		} else {
			panic(exception.NewValidationError("Bundle item tidak boleh kosong"))
		}
	}
	// Product Variant

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

	res, err := service.ProductRepository.GetById(tx, productRes.Id)
	exception.PanicIfNeeded(err)

	return service.MapProduct(res, currentUserId != 0)
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
	product.IsShowOnlyInMarketPlace = payload.IsShowOnlyInMarketPlace
	product.ImageIds = helpers.JoinImageIds(payload.ImageIds)
	product.VideoUrl = payload.VideoUrl
	product.ComponentTypeId = payload.ComponentTypeId
	product.Tags = payload.Tags
	product.ProductCategoryId = payload.ProductCategoryId
	product.BrandId = payload.BrandId
	product.Description = payload.Description
	product.ProductSpec = payload.ProductSpec
	product.ModifiedById = currentUserId
	product.ModifiedAt = time.Now().UTC()

	_, err = service.ProductRepository.Save(tx, product)
	exception.PanicIfNeeded(err)

	// Product SKU (Update, Delete, Add New)
	productSkus := product.ProductSkus
	productSkuToBeDeleted := []entity.ProductSku{}
	productGroupItems := []entity.ProductGroupItem{}
	productGroupItemsToBeDeleted := []entity.ProductGroupItem{}
	componentSpecs := []entity.ComponentSpec{}
	componentSpecsToBeDeleted := []entity.ComponentSpec{}

	// Update, delete
	for i, productSku := range productSkus {
		payloadIdx := slices.IndexFunc(payload.ProductSkus, func(model request.ProductSkuPayload) bool {
			return model.Id == productSku.Id
		})

		if payloadIdx != -1 {
			productSkus[i].SKU = payload.ProductSkus[payloadIdx].SKU
			productSkus[i].Name = payload.ProductSkus[payloadIdx].Name
			productSkus[i].Price = payload.ProductSkus[payloadIdx].Price
			productSkus[i].Stock = payload.ProductSkus[payloadIdx].Stock
			productSkus[i].StockAlert = payload.ProductSkus[payloadIdx].StockAlert
			productSkus[i].IsActive = payload.ProductSkus[payloadIdx].IsActive

			productSkuSpecs := productSku.ComponentSpecs

			// update exsisting componentSpec
			for i, componentSpec := range productSkuSpecs {
				componentSpecIdx := slices.IndexFunc(payload.ProductSkus[payloadIdx].ComponentSpecs, func(model request.ComponentSpec) bool {
					return model.Id == componentSpec.Id
				})

				if componentSpecIdx != -1 {
					componentSpecPayload := payload.ProductSkus[payloadIdx].ComponentSpecs[componentSpecIdx]
					productSkuSpecs[i].SpecKey = componentSpecPayload.SpecKey
					productSkuSpecs[i].SpecValue = componentSpecPayload.SpecValue

				} else {
					componentSpecsToBeDeleted = append(componentSpecsToBeDeleted, componentSpec)
				}
			}

			// add new component spec
			for _, productSkuSpec := range payload.ProductSkus[payloadIdx].ComponentSpecs {

				savedIdx := slices.IndexFunc(productSkuSpecs, func(model entity.ComponentSpec) bool {
					return model.Id == productSkuSpec.Id
				})

				if savedIdx == -1 {
					productSkuSpecs = append(productSkuSpecs, entity.ComponentSpec{
						Id:           productSkuSpec.Id,
						ProductSkuId: productSku.Id,
						SpecKey:      productSkuSpec.SpecKey,
						SpecValue:    productSkuSpec.SpecValue,
						Sequence:     productSkuSpecs[len(productSkuSpecs)-1].Sequence + 1,
					})
				}
			}

			componentSpecs = append(componentSpecs, productSkuSpecs...)

			// Product Bundle items (Update, Delete, Add New)
			if product.ProductType == enum.ProductTypeGroup {
				if len(payload.ProductSkus[payloadIdx].ProductGroupItems) == 0 {
					panic(exception.NewValidationError("Group detail tidak boleh kosong"))
				}

				productSkuGroupItems := productSku.ProductGroupItems
				for i, groupItem := range productSkuGroupItems {
					groupItemIdx := slices.IndexFunc(payload.ProductSkus[payloadIdx].ProductGroupItems, func(model request.ProductGroupItemPayload) bool {
						return model.Id == groupItem.Id
					})

					if groupItemIdx != -1 {
						productGroupPayload := payload.ProductSkus[payloadIdx].ProductGroupItems[groupItemIdx]

						productSkuGroupItems[i].ProductId = productGroupPayload.ProductId
						productSkuGroupItems[i].ProductSkuId = productGroupPayload.ProductSkuId
						productSkuGroupItems[i].Qty = productGroupPayload.Qty
					} else {
						productGroupItemsToBeDeleted = append(productGroupItemsToBeDeleted, groupItem)
					}
				}

				// Add new
				for _, groupItem := range payload.ProductSkus[payloadIdx].ProductGroupItems {
					savedIdx := slices.IndexFunc(productSku.ProductGroupItems, func(model entity.ProductGroupItem) bool {
						return model.Id == groupItem.Id
					})

					if savedIdx == -1 {
						productSkuGroupItems = append(productSkuGroupItems, entity.ProductGroupItem{
							Id:           groupItem.Id,
							ParentId:     productSku.Id,
							ProductId:    groupItem.ProductId,
							ProductSkuId: groupItem.ProductSkuId,
							Qty:          groupItem.Qty,
							Sequence:     productSkuGroupItems[len(productSkuGroupItems)-1].Sequence + 1,
						})
					}
				}

				productGroupItems = append(productGroupItems, productSkuGroupItems...)
			}
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
				Id:         productSku.Id,
				ProductId:  product.Id,
				SKU:        productSku.SKU,
				Price:      productSku.Price,
				Stock:      productSku.Stock,
				StockAlert: productSku.StockAlert,
				IsActive:   productSku.IsActive,
				Sequence:   productSkus[len(productSkus)-1].Sequence + 1,
			})

			// Product Spec
			for _, componentSpec := range productSku.ComponentSpecs {
				newSequence := 1
				if len(componentSpecs) > 0 {
					newSequence = componentSpecs[len(componentSpecs)-1].Sequence + 1
				}
				componentSpecs = append(componentSpecs, entity.ComponentSpec{
					Id:           componentSpec.Id,
					ProductSkuId: productSku.Id,
					SpecKey:      componentSpec.SpecKey,
					SpecValue:    componentSpec.SpecValue,
					Sequence:     newSequence,
				})
			}

			// Product Group
			for _, groupItem := range productSku.ProductGroupItems {
				newSequence := 1
				if len(productGroupItems) > 0 {
					newSequence = productGroupItems[len(productGroupItems)-1].Sequence + 1
				}
				productGroupItems = append(productGroupItems, entity.ProductGroupItem{
					Id:           groupItem.Id,
					ParentId:     productSku.Id,
					ProductId:    groupItem.ProductId,
					ProductSkuId: groupItem.ProductSkuId,
					Qty:          groupItem.Qty,
					Sequence:     newSequence,
				})
			}
		}

	}

	err = service.ProductRepository.SaveProductSkus(tx, productSkus)
	exception.PanicIfNeeded(err)

	err = service.ProductRepository.DeleteProductSkus(tx, product.Id, productSkuToBeDeleted)
	exception.PanicIfNeeded(err)

	err = service.ProductRepository.SaveProductGroupItems(tx, productGroupItems)
	exception.PanicIfNeeded(err)

	err = service.ProductRepository.DeleteProductGroupItems(tx, productGroupItemsToBeDeleted)
	exception.PanicIfNeeded(err)

	err = service.ProductRepository.SaveComponentSpecs(tx, componentSpecs)
	exception.PanicIfNeeded(err)

	err = service.ProductRepository.DeleteComponentSpecs(tx, componentSpecsToBeDeleted)
	exception.PanicIfNeeded(err)

	// Product Variant
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
			newSequence := 1
			if len(productVariantOptions) > 0 {
				newSequence = productVariantOptions[len(productVariantOptions)-1].Sequence + 1
			}
			productVariantOptions = append(productVariantOptions, entity.ProductVariantOption{
				Id:        variantOption.Id,
				ProductId: product.Id,
				Name:      variantOption.Name,
				Sequence:  newSequence,
			})
		}
	}

	err = service.ProductRepository.SaveProductVariantOptions(tx, productVariantOptions)
	exception.PanicIfNeeded(err)

	err = service.ProductRepository.DeleteProductVariantOptions(tx, productVariantOptionsToBeDeleted)
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
			newSequence := 1
			if len(productVariantOptionValues) > 0 {
				newSequence = productVariantOptionValues[len(productVariantOptionValues)-1].Sequence + 1
			}
			productVariantOptionValues = append(productVariantOptionValues, entity.ProductVariantOptionValue{
				Id:                     variantOptionValue.Id,
				ProductId:              product.Id,
				ProductVariantOptionId: variantOptionValue.ProductVariantOptionId,
				Name:                   variantOptionValue.Name,
				Sequence:               newSequence,
			})
		}
	}

	err = service.ProductRepository.SaveProductVariantOptionValues(tx, productVariantOptionValues)
	exception.PanicIfNeeded(err)

	err = service.ProductRepository.DeleteProductVariantOptionValues(tx, productVariantOptionValuesToBeDeleted)
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
			newSequence := 1
			if len(productVariantOptionValues) > 0 {
				newSequence = productVariantOptionValues[len(productVariantOptionValues)-1].Sequence + 1
			}
			productSkuVariants = append(productSkuVariants, entity.ProductSkuVariant{
				Id:                          skuVariant.Id,
				ProductId:                   product.Id,
				ProductSkuId:                skuVariant.ProductSkuId,
				ProductVariantOptionId:      skuVariant.ProductVariantOptionId,
				ProductVariantOptionValueId: skuVariant.ProductVariantOptionValueId,
				Sequence:                    newSequence,
			})
		}
	}

	err = service.ProductRepository.SaveProductSkuVariants(tx, productSkuVariants)
	exception.PanicIfNeeded(err)

	err = service.ProductRepository.DeleteProductSkuVariants(tx, productSkuVariantsToBeDeleted)
	exception.PanicIfNeeded(err)

	// return new result
	res, err := service.ProductRepository.GetById(tx, productId)
	exception.PanicIfNeeded(err)

	return service.MapProduct(res, currentUserId != 0)
}

func (service *ProductServiceImpl) ChangeComponent(currentUserId uint, payload request.ChangeComponent) string {
	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	// Get Data
	oldProduct := entity.Product{}
	err := tx.Model(oldProduct).
		Joins("Join product_skus on product_skus.product_id = products.id and product_skus.id = ?", payload.OldProductSkuId).
		Preload("ProductSkus").First(&oldProduct).Error

	exception.PanicIfNeeded(err)

	newProduct := entity.Product{}
	err = tx.Model(newProduct).
		Joins("Join product_skus on product_skus.product_id = products.id and product_skus.id = ?", payload.NewProductSkuId).
		Preload("ProductSkus").First(&newProduct).Error

	exception.PanicIfNeeded(err)

	// validate new product

	if !newProduct.IsActive || !newProduct.ProductSkus[0].IsActive {
		productName := newProduct.Name
		if newProduct.ProductSkus[0].Name != "" {
			productName += " / " + newProduct.ProductSkus[0].Name
		}
		panic(exception.NewValidationError(fmt.Sprintf("Produk %s tidak aktif", productName)))
	}

	res := tx.Exec("UPDATE product_group_items set product_id = ?, product_sku_id = ? where product_sku_id = ?",
		newProduct.Id, newProduct.ProductSkus[0].Id, payload.OldProductSkuId)

	exception.PanicIfNeeded(res.Error)

	return fmt.Sprintf("%s produk telah diperbarui", strconv.Itoa(int(res.RowsAffected)))
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

func (service *ProductServiceImpl) Search(currentUserId uint, query *string, ids *[]uuid.UUID, productTypes *[]string, componentTypeIds *[]uint, productCategoryIds *[]uint, brandIds *[]uint, tag *string, minPrice, maxPrice *float64, isInStockOnly *bool, isActive, isShowOnlyInMarketPlace *bool, page, pageSize *int) response.PageResult {

	res, totalCount, totalPage := service.ProductRepository.Search(service.db, currentUserId != 0, query, ids, productTypes, componentTypeIds, productCategoryIds, brandIds, tag, minPrice, maxPrice, isInStockOnly, isActive, isShowOnlyInMarketPlace, page, pageSize)

	return response.PageResult{
		Items:      service.mapProducts(res, currentUserId != 0),
		TotalCount: totalCount,
		PageSize:   totalPage,
	}
}

func (service *ProductServiceImpl) BrowseProductSku(currentUserId uint, query *string, productTypes *[]string, productCategoryId *uint, page, pageSize *int) response.PageResult {
	res, count, totalPage := service.ProductRepository.BrowseProductSku(service.db, query, productTypes, productCategoryId, page, pageSize)
	return response.PageResult{
		Items:      res,
		TotalCount: count,
		PageSize:   totalPage,
	}
}

func (service *ProductServiceImpl) GetById(currentUserId uint, productId uuid.UUID) response.ProductResponse {
	res, err := service.ProductRepository.GetById(service.db, productId)
	exception.PanicIfNeeded(err)

	return service.MapProduct(res, currentUserId != 0)
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
func (service *ProductServiceImpl) mapProducts(products []entity.Product, isAdmin bool) []response.ProductResponse {
	productRes := []response.ProductResponse{}

	for _, el := range products {
		productRes = append(productRes, service.MapProduct(el, isAdmin))
	}

	return productRes
}

func (service *ProductServiceImpl) MapProduct(product entity.Product, isAdmin bool) response.ProductResponse {
	productRes := response.ProductResponse{
		Id:                      product.Id,
		SKU:                     product.SKU,
		Name:                    product.Name,
		IsActive:                product.IsActive,
		IsShowOnlyInMarketPlace: product.IsShowOnlyInMarketPlace,
		ImageIds:                helpers.SplitImageIds(product.ImageIds),
		VideoUrl:                product.VideoUrl,
		Tags:                    product.Tags,
		Description:             product.Description,
		ProductSpec:             product.ProductSpec,
		CreatedBy:               data_mapper.MapAuditTrail(product.CreatedBy),
		CreatedAt:               product.CreatedAt,
		ModifiedBy:              data_mapper.MapAuditTrail(product.ModifiedBy),
		ModifiedAt:              product.ModifiedAt,

		ProductSkus:                []response.ProductSku{},
		ProductVariantOptions:      []response.ProductVariantOption{},
		ProductVariantOptionValues: []response.ProductVariantOptionValue{},
		ProductSkuVariants:         []response.ProductSkuVariant{},
	}

	productRes.ProductType = response.EnumResponse{
		Code: string(product.ProductType),
		Name: product.ProductType.DisplayString(),
	}

	if product.ProductCategory != nil {
		productCategoryRes := response.ProductCategory{
			Id:   product.ProductCategory.Id,
			Name: product.ProductCategory.Name,
		}
		productRes.ProductCategory = &productCategoryRes
	}

	if product.Brand != nil {
		brandRes := response.Brand{
			Id:   product.Brand.Id,
			Name: product.Brand.Name,
		}
		productRes.Brand = &brandRes
	}

	if product.ComponentType != nil {
		componentTypeRes := response.ComponentType{
			Id:   product.ComponentType.Id,
			Code: product.ComponentType.Code,
			Name: product.ComponentType.Name,
		}
		productRes.ComponentType = &componentTypeRes
	}

	for _, productSku := range product.ProductSkus {
		if !isAdmin && !productSku.IsActive {
			continue
		}
		productSkuRes := response.ProductSku{
			Id:                productSku.Id,
			ProductId:         productSku.ProductId,
			SKU:               productSku.SKU,
			Name:              productSku.Name,
			Price:             productSku.Price,
			Stock:             productSku.Stock,
			StockAlert:        productSku.StockAlert,
			IsActive:          productSku.IsActive,
			Sequence:          productSku.Sequence,
			ComponentSpecs:    service.mapComponentSpecs(productSku.ComponentSpecs),
			ProductGroupItems: service.mapProductGroupItems(productSku.ProductGroupItems),
		}

		productRes.ProductSkus = append(productRes.ProductSkus, productSkuRes)
	}

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

	return productRes
}

func (service *ProductServiceImpl) mapComponentSpecs(componentSpecs []entity.ComponentSpec) []response.ComponentSpec {

	componentSpecRes := []response.ComponentSpec{}
	for _, componentSpec := range componentSpecs {
		componentSpecRes = append(componentSpecRes, response.ComponentSpec{
			Id:           componentSpec.Id,
			ProductSkuId: componentSpec.ProductSkuId,
			SpecKey:      componentSpec.SpecKey,
			SpecValue:    componentSpec.SpecValue,
		})
	}

	return componentSpecRes
}

func (service *ProductServiceImpl) mapProductGroupItems(productGroupItems []entity.ProductGroupItem) []response.ProductGroupItemResponse {
	res := []response.ProductGroupItemResponse{}

	for _, item := range productGroupItems {
		res = append(res, response.ProductGroupItemResponse{
			Id:           item.Id,
			ProductId:    item.Product.Id,
			ProductSkuId: item.ProductSkuId,
			Qty:          item.Qty,
			Product:      service.MapProduct(item.Product, true),
		})
	}

	return res
}

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
