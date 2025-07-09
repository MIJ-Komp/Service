package service

import (
	"api.mijkomp.com/models/entity"
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
	"github.com/google/uuid"
)

type ProductService interface {
	Create(currentUserId uint, productId uuid.UUID, req request.ProductPayload) response.ProductResponse
	Update(currentUserId uint, productId uuid.UUID, req request.ProductPayload) response.ProductResponse
	Delete(currentUserId uint, productId uuid.UUID) string
	Search(currentUserId uint, query *string, ids *[]uuid.UUID, productTypes *[]string, productCategoryIds *[]uint, componentTypeIds *[]uint, isActive, isShowOnlyInMarketPlace *bool, page, pageSize *int) response.PageResult
	BrowseProductSku(currentUserId uint, query *string, productTypes *[]string, productCategoryId *uint, page, pageSize *int) response.PageResult
	GetById(currentUserId uint, productId uuid.UUID) response.ProductResponse

	// Variant Option
	CreateVariantOptions(currentUserId uint, payload request.VariantOptionPayload) response.VariantOption
	GetVariantOptions(currentUserId uint) []response.VariantOption
	MapProduct(product entity.Product, isAdmin bool) response.ProductResponse
}
