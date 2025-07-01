package service_impl

import (
	"time"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers"
	"api.mijkomp.com/helpers/data_mapper"
	"api.mijkomp.com/models/entity"
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/repository"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type MenuServiceImpl struct {
	MenuRepository repository.MenuRepository
	db             *gorm.DB
	Validation     *validator.Validate
}

func NewMenuService(menuRepostitory repository.MenuRepository, validation *validator.Validate, db *gorm.DB) *MenuServiceImpl {
	return &MenuServiceImpl{
		MenuRepository: menuRepostitory,
		Validation:     validation,
		db:             db,
	}
}

func (service *MenuServiceImpl) Create(currentUserId uint, payload request.Menu) response.Menu {

	err := service.Validation.Struct(payload)
	exception.PanicIfNeeded(err)

	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	menuEntity := entity.Menu{
		Name:         payload.Name,
		ParentId:     payload.ParentId,
		Path:         payload.Path,
		CreatedById:  currentUserId,
		CreatedAt:    time.Now(),
		ModifiedById: currentUserId,
		ModifiedAt:   time.Now(),
	}

	result, err := service.MenuRepository.Save(tx, menuEntity)
	exception.PanicIfNeeded(err)

	return service.GenerateGetResult(result)
}

func (service *MenuServiceImpl) Update(currentUserId uint, menuId uint, payload request.Menu) response.Menu {

	err := service.Validation.Struct(payload)
	exception.PanicIfNeeded(err)

	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	menu, err := service.MenuRepository.GetById(tx, menuId)
	exception.PanicIfNeeded(err)

	menu.Name = payload.Name
	menu.ParentId = payload.ParentId
	menu.Path = payload.Path
	menu.ModifiedById = currentUserId
	menu.ModifiedAt = time.Now()

	result, err := service.MenuRepository.Save(tx, menu)
	exception.PanicIfNeeded(err)

	return service.GenerateGetResult(result)
}

func (service *MenuServiceImpl) Delete(currentUserId uint, menuId uint) string {
	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	err := service.MenuRepository.Delete(tx, menuId)
	exception.PanicIfNeeded(err)

	return "Kategori berhasil dihapus"
}

func (service *MenuServiceImpl) Search(currentUserId uint, query *string, parentId *uint) []response.Menu {
	res := service.MenuRepository.Search(service.db, query, parentId)

	return service.GenerateSearchResult(res)
}

// func (service *MenuServiceImpl) GetById(currentUserId uint, menuId uint) response.Menu {
// 	res, err := service.MenuRepository.GetById(service.db, menuId)
// 	exception.PanicIfNeeded(err)

// 	return service.GenerateGetResult(res)
// }

func (service *MenuServiceImpl) CreateItem(currentUserId uint, menuId uint, payload request.MenuItem) string {

	err := service.Validation.Struct(payload)
	exception.PanicIfNeeded(err)

	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	menuItemEntity := entity.MenuItem{
		MenuId:            menuId,
		ProductCategoryId: payload.ProductCategoryId,
	}

	err = service.MenuRepository.CreateItem(tx, menuItemEntity)
	exception.PanicIfNeeded(err)

	return "Menu item berhasil ditambah"
}

func (service *MenuServiceImpl) DeleteItem(currentUserId uint, itemId uint) string {
	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	err := service.MenuRepository.DeleteItem(tx, itemId)
	exception.PanicIfNeeded(err)

	return "Menu item berhasil dihapus"
}

// map helpers
func (service *MenuServiceImpl) GenerateSearchResult(menus []entity.Menu) []response.Menu {

	res := []response.Menu{}
	for _, menu := range menus {

		menuRes := response.Menu{
			Id:         menu.Id,
			Name:       menu.Name,
			ParentId:   menu.ParentId,
			Path:       menu.Path,
			CreatedBy:  data_mapper.MapAuditTrail(menu.CreatedBy),
			CreatedAt:  menu.CreatedAt,
			ModifiedBy: data_mapper.MapAuditTrail(menu.ModifiedBy),
			ModifiedAt: menu.ModifiedAt,

			MenuItems: service.GenerateMenuItemResult(menu.MenuItems),
			Childs:    service.GenerateChildResult(menu.Id, menus),
		}

		if menu.ParentId == nil {
			res = append(res, menuRes)
		}

	}

	return res
}

func (service *MenuServiceImpl) GenerateChildResult(parentId uint, menus []entity.Menu) []response.Menu {
	childs := []response.Menu{}

	for _, menuChild := range menus {
		if menuChild.ParentId != nil && parentId == *menuChild.ParentId {
			childs = append(childs, response.Menu{
				Id:         menuChild.Id,
				Name:       menuChild.Name,
				ParentId:   menuChild.ParentId,
				Path:       menuChild.Path,
				CreatedBy:  data_mapper.MapAuditTrail(menuChild.CreatedBy),
				CreatedAt:  menuChild.CreatedAt,
				ModifiedBy: data_mapper.MapAuditTrail(menuChild.ModifiedBy),
				ModifiedAt: menuChild.ModifiedAt,

				MenuItems: service.GenerateMenuItemResult(menuChild.MenuItems),
				Childs:    service.GenerateChildResult(menuChild.Id, menus),
			})
		}

	}
	return childs
}

func (service *MenuServiceImpl) GenerateGetResult(menu entity.Menu) response.Menu {
	res := response.Menu{
		Id:         menu.Id,
		Name:       menu.Name,
		ParentId:   menu.ParentId,
		Path:       menu.Path,
		CreatedBy:  data_mapper.MapAuditTrail(menu.CreatedBy),
		CreatedAt:  menu.CreatedAt,
		ModifiedBy: data_mapper.MapAuditTrail(menu.ModifiedBy),
		ModifiedAt: menu.ModifiedAt,

		MenuItems: service.GenerateMenuItemResult(menu.MenuItems),
	}

	return res
}

func (service *MenuServiceImpl) GenerateMenuItemResult(menuItems []entity.MenuItem) []response.MenuItem {
	result := []response.MenuItem{}

	for _, item := range menuItems {
		result = append(result, response.MenuItem{
			Id:                item.Id,
			ProductCategoryId: item.ProductCategoryId,
			Name:              item.ProductCategory.Name,
		})
	}
	return result
}
