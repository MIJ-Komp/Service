package customer

import (
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
	menu := app.Group("/api/menus")
	menu.Get("/", controller.Search)
}

// @Summary		Search menu
// @Tags			Menu
// @Accept		json
// @Produce		json
// @Success		200	{object}	response.Menu
// @Failure		404	{object}	response.WebResponse
// @Router		/api/menus [get]
func (controller *MenuController) Search(ctx *fiber.Ctx) error {
	result := controller.MenuService.Search(0, nil, nil)

	return ctx.JSON(response.NewWebResponse(result))
}
