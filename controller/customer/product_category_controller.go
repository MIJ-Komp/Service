package customer

import (
	"api.mijkomp.com/helpers"
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
	category := app.Group("/api/product-categories")
	category.Get("/", controller.Search)
	category.Get("/:id", controller.GetById)
}

// @Summary		Search product categories
// @Tags			Product Category
// @Accept		json
// @Produce		json
// @Param			query query string false " "
// @Param			parentId query int false " "
// @Success		200	{object} response.ProductCategory
// @Failure		404	{object} response.WebResponse
// @Router		/api/product-categories [get]
func (controller *ProductCategoryController) Search(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	query := ctx.Query("query")
	parentId := helpers.ParseNullableUint(ctx.Query("parentId"))

	result := controller.ProductCategoryService.Search(currentUserId, &query, parentId)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary		Get category by id
// @Tags			Product Category
// @Accept		json
// @Produce		json
// @Param			id path int true " "
// @Success		200	{object} response.ProductCategory
// @Failure		404	{object} response.WebResponse
// @Router		/api/product-categories/{id} [get]
func (controller *ProductCategoryController) GetById(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	id := helpers.ParseUint(ctx.Params("id"))

	result := controller.ProductCategoryService.GetById(currentUserId, id)

	return ctx.JSON(response.NewWebResponse(result))
}
