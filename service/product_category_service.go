package service

import (
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
)

type ProductCategoryService interface {
	Create(currentUserId uint, payload request.ProductCategory) response.ProductCategory
	Update(currentUserId uint, categoryId uint, payload request.ProductCategory) response.ProductCategory
	Delete(currentUserId uint, categoryId uint) string
	Search(currentUserId uint, query *string) []response.ProductCategory
	GetById(currentUserId uint, categoryId uint) response.ProductCategory
}
