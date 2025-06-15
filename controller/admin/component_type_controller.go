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

type ComponentTypeController struct {
	UserService          service.UserService
	ComponentTypeService service.ComponentTypeService
}

func NewComponentTypeController(userService *service.UserService, componentTypeService *service.ComponentTypeService) *ComponentTypeController {
	return &ComponentTypeController{UserService: *userService, ComponentTypeService: *componentTypeService}
}

func (controller *ComponentTypeController) Route(app *fiber.App) {
	componentType := app.Group("/api/admin/component-types", middleware.AuthMiddleware(controller.UserService))
	componentType.Post("/", controller.Create)
	componentType.Put("/:id", controller.Update)
	componentType.Delete("/:id", controller.Delete)
	componentType.Get("/", controller.Search)
	componentType.Get("/:id", controller.GetById)
}

// @Summary      Create Component ype
// @Tags         Component Type
// @Accept       json
// @Produce      json
// @Param        request body request.ComponentType  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security		 ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/component-types [post]
func (controller *ComponentTypeController) Create(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	var payload request.ComponentType
	err := ctx.BodyParser(&payload)
	exception.PanicIfNeeded(err)

	result := controller.ComponentTypeService.Create(currentUserId, payload)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary      Update component type
// @Tags         Component Type
// @Accept       json
// @Produce      json
// @Param        id path int  true " "
// @Param        request body request.ComponentType  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/component-types/{id} [put]
func (controller *ComponentTypeController) Update(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))
	componentTypeId := helpers.ParseUserId(ctx.Params("id"))

	var payload request.ComponentType
	err := ctx.BodyParser(&payload)
	exception.PanicIfNeeded(err)

	result := controller.ComponentTypeService.Update(currentUserId, componentTypeId, payload)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary      Delete component type
// @Tags         Component Type
// @Accept       json
// @Produce      json
// @Param        id path int  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security		 ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/component-types/{id} [delete]
func (controller *ComponentTypeController) Delete(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))
	componentTypeId := helpers.ParseUserId(ctx.Params("id"))

	result := controller.ComponentTypeService.Delete(currentUserId, componentTypeId)

	return ctx.JSON(response.NewWebResponse(nil, result))
}

// @Summary		Search component type
// @Tags			Component Type
// @Accept		json
// @Produce		json
// @Param			query query string false " "
// @Success		200	{object}	response.ComponentType
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/component-types [get]
func (controller *ComponentTypeController) Search(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	query := ctx.Query("query")

	result := controller.ComponentTypeService.Search(currentUserId, &query)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary		Get component type by id
// @Tags			Component Type
// @Accept		json
// @Produce		json
// @Param			id path int true " "
// @Success		200	{object}	response.ComponentType
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/component-types/{id} [get]
func (controller *ComponentTypeController) GetById(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	id := helpers.ParseUint(ctx.Params("id"))

	result := controller.ComponentTypeService.GetById(currentUserId, id)

	return ctx.JSON(response.NewWebResponse(result))
}
