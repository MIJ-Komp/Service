package service

import (
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
)

type BrandService interface {
	Create(currentUserId uint, payload request.Brand) response.Brand
	Update(currentUserId uint, brandId uint, payload request.Brand) response.Brand
	Delete(currentUserId uint, brandId uint) string
	Search(currentUserId uint, query *string) []response.Brand
	GetById(currentUserId uint, brandId uint) response.Brand
}
