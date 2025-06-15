package service

import (
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
)

type ComponentTypeService interface {
	Create(currentUserId uint, payload request.ComponentType) response.ComponentType
	Update(currentUserId uint, componentTypeId uint, payload request.ComponentType) response.ComponentType
	Delete(currentUserId uint, componentTypeId uint) string
	Search(currentUserId uint, query *string) []response.ComponentType
	GetById(currentUserId uint, componentTypeId uint) response.ComponentType
}
