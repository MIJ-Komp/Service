package admin

import (
	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers"
	"api.mijkomp.com/middleware"
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/service"
	"github.com/gofiber/fiber/v2"
)

type MenuController struct {
	UserService service.UserService
	MenuService service.MenuService
}

func NewMenuController(userService *service.UserService, menuService *service.MenuService) *MenuController {
	return &MenuController{UserService: *userService, MenuService: *menuService}
}

func (controller *MenuController) Route(app *fiber.App) {
	menu := app.Group("/api/admin/menus", middleware.AuthMiddleware(controller.UserService))
	menu.Post("/", controller.Create)
	menu.Put("/:id", controller.Update)
	menu.Delete("/:id", controller.Delete)
	menu.Get("/", controller.Search)
	// menu.Get("/:id", controller.GetById)
	menu.Post("/:menuId", controller.CreateItem)
	menu.Delete("/delete-item/:id", controller.DeleteItem)
}

// @Summary      Create menu
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Param        request body request.Menu  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/menus [post]
func (controller *MenuController) Create(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	var payload request.Menu
	err := ctx.BodyParser(&payload)
	exception.PanicIfNeeded(err)

	result := controller.MenuService.Create(currentUserId, payload)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary      Update menu
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Param        id path int  true " "
// @Param        request body request.Menu  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/menus/{id} [put]
func (controller *MenuController) Update(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))
	menuId := helpers.ParseUserId(ctx.Params("id"))

	var payload request.Menu
	err := ctx.BodyParser(&payload)
	exception.PanicIfNeeded(err)

	result := controller.MenuService.Update(currentUserId, menuId, payload)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary      Delete menu
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Param        id path int  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/menus/{id} [delete]
func (controller *MenuController) Delete(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))
	menuId := helpers.ParseUserId(ctx.Params("id"))

	result := controller.MenuService.Delete(currentUserId, menuId)

	return ctx.JSON(response.NewWebResponse(nil, result))
}

// @Summary		Search menu
// @Tags			Menu
// @Accept		json
// @Produce		json
// @Param			query query string false " "
// @Param			parentId query int false " "
// @Success		200	{object}	response.Menu
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/menus [get]
func (controller *MenuController) Search(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	query := ctx.Query("query")
	parentId := helpers.ParseNullableUint(ctx.Query("parentId"))

	result := controller.MenuService.Search(currentUserId, &query, parentId)

	return ctx.JSON(response.NewWebResponse(result))
}

// // @Summary		Get menu by id
// // @Tags			Menu
// // @Accept		json
// // @Produce		json
// // @Param			id path int true " "
// // @Success		200	{object}	response.Menu
// // @Failure		404	{object}	response.WebResponse
// // @Security	ApiKeyAuth
// // @in header
// // @name Authorization
// // @Router		/api/admin/menus/{id} [get]
// func (controller *MenuController) GetById(ctx *fiber.Ctx) error {
// 	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

// 	id := helpers.ParseUint(ctx.Params("id"))

// 	result := controller.MenuService.GetById(currentUserId, id)

// 	return ctx.JSON(response.NewWebResponse(result))
// }

// @Summary      Create menu item
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Param				 menuId path int true " "
// @Param        request body request.MenuItem  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/menus/{menuId} [post]
func (controller *MenuController) CreateItem(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	id := helpers.ParseUint(ctx.Params("menuId"))

	var payload request.MenuItem
	err := ctx.BodyParser(&payload)
	exception.PanicIfNeeded(err)

	result := controller.MenuService.CreateItem(currentUserId, id, payload)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary      Delete menu item
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Param			 	 id path int true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/menus/delete-item/{id} [delete]
func (controller *MenuController) DeleteItem(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	// id := helpers.ParseUint(ctx.Params("id"))
	id := helpers.ParseUint(ctx.Params("id"))

	result := controller.MenuService.DeleteItem(currentUserId, id)

	return ctx.JSON(response.NewWebResponse(result))
}
