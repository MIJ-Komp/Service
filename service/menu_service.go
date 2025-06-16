package service

import (
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
)

type MenuService interface {
	Create(currentUserId uint, payload request.Menu) response.Menu
	Update(currentUserId uint, menuId uint, payload request.Menu) response.Menu
	Delete(currentUserId uint, menuId uint) string
	Search(currentUserId uint, query *string, parentId *uint) []response.Menu
	// GetById(currentUserId uint, menuId uint) response.Menu

	CreateItem(currentUserId uint, menuId uint, payload request.MenuItem) string
	DeleteItem(currentUserId uint, itemId uint) string
}
