package controller

import (
	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers"
	"api.mijkomp.com/middleware"
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/service"
	"github.com/gofiber/fiber/v2"
)

type ProductCategoryController struct {
	UserService            service.UserService
	ProductCategoryService service.ProductCategoryService
}

func NewProductCategoryController(userService *service.UserService, categoryService *service.ProductCategoryService) *ProductCategoryController {
	return &ProductCategoryController{UserService: *userService, ProductCategoryService: *categoryService}
}

func (controller *ProductCategoryController) Route(app *fiber.App) {
	category := app.Group("/api/admin/product-categories", middleware.AuthMiddleware(controller.UserService))
	category.Post("/", controller.Create)
	category.Put("/:id", controller.Update)
	category.Delete("/:id", controller.Delete)
	category.Get("/", controller.Search)
	category.Get("/:id", controller.GetById)
}

// @Summary      Create category
// @Tags         (Admin) Product Category
// @Accept       json
// @Produce      json
// @Param        request body request.ProductCategory  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/product-categories [post]
func (controller *ProductCategoryController) Create(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	var payload request.ProductCategory
	err := ctx.BodyParser(&payload)
	exception.PanicIfNeeded(err)

	result := controller.ProductCategoryService.Create(currentUserId, payload)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary      Update category
// @Tags         (Admin) Product Category
// @Accept       json
// @Produce      json
// @Param        id path int  true " "
// @Param        request body request.ProductCategory  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/product-categories/{id} [put]
func (controller *ProductCategoryController) Update(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))
	categoryId := helpers.ParseUserId(ctx.Params("id"))

	var payload request.ProductCategory
	err := ctx.BodyParser(&payload)
	exception.PanicIfNeeded(err)

	result := controller.ProductCategoryService.Update(currentUserId, categoryId, payload)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary      Delete category
// @Tags         (Admin) Product Category
// @Accept       json
// @Produce      json
// @Param        id path int  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/product-categories/{id} [delete]
func (controller *ProductCategoryController) Delete(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))
	categoryId := helpers.ParseUserId(ctx.Params("id"))

	result := controller.ProductCategoryService.Delete(currentUserId, categoryId)

	return ctx.JSON(response.NewWebResponse(nil, result))
}

// @Summary		Search category
// @Tags			(Admin) Product Category
// @Accept		json
// @Produce		json
// @Param			query query string false " "
// @Param			parentId query int false " "
// @Success		200	{object}	response.ProductCategory
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/product-categories [get]
func (controller *ProductCategoryController) Search(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	query := ctx.Query("query")
	parentId := helpers.ParseNullableUint(ctx.Query("parentId"))

	result := controller.ProductCategoryService.Search(currentUserId, &query, parentId)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary		Get category by id
// @Tags			(Admin) Product Category
// @Accept		json
// @Produce		json
// @Param			id path int true " "
// @Success		200	{object}	response.ProductCategory
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/product-categories/{id} [get]
func (controller *ProductCategoryController) GetById(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	id := helpers.ParseUint(ctx.Params("id"))

	result := controller.ProductCategoryService.GetById(currentUserId, id)

	return ctx.JSON(response.NewWebResponse(result))
}
