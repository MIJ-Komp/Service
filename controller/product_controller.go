package controller

import (
	"strings"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers"
	"api.mijkomp.com/middleware"
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/service"
	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	ProductService service.ProductService
	UserService    service.UserService
}

func NewProductController(userService *service.UserService, productService *service.ProductService) *ProductController {
	return &ProductController{UserService: *userService, ProductService: *productService}
}

func (controller *ProductController) Route(app *fiber.App) {

	product := app.Group("/api/admin/products")
	product.Post("/:id", middleware.AuthMiddleware(controller.UserService), controller.Create)
	product.Put("/:id", middleware.AuthMiddleware(controller.UserService), controller.Update)
	product.Delete("/:id", middleware.AuthMiddleware(controller.UserService), controller.Delete)
	product.Get("/", middleware.AuthMiddleware(controller.UserService), controller.Search)
	// product.Get("/browse", middleware.AuthMiddleware(controller.UserService), controller.BrowseProductSku)
	product.Get("/:id", middleware.AuthMiddleware(controller.UserService), controller.GetById)

	variantOption := app.Group("/api/variant-options")
	variantOption.Post("/", middleware.AuthMiddleware(controller.UserService), controller.CreateVariantOptions)
	variantOption.Get("/", middleware.AuthMiddleware(controller.UserService), controller.GetVariantOptions)
}

// @Summary     Create Product
// @Tags        (Admin) Product
// @Accept      json
// @Produce     json
// @Param       id path string true "Id (UUID format)" format(uuid)
// @Param       request body request.ProductPayload true " "
// @Success     200  {object}  response.WebResponse
// @Failure     400  {object}  response.WebResponse
// @Security		ApiKeyAuth
// @in 					header
// @name 				Authorization
// @Router       /api/admin/products/{id} [post]
func (controller *ProductController) Create(ctx *fiber.Ctx) error {

	currentUserId := helpers.ParseUint(ctx.Locals("userId").(string))
	productId := helpers.ParseUUID(ctx.Params("id"))

	var productModel request.ProductPayload
	err := ctx.BodyParser(&productModel)
	exception.NewModelValidationError(err)

	result := controller.ProductService.Create(currentUserId, productId, productModel)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary     Update Product
// @Tags        (Admin) Product
// @Accept      json
// @Produce     json
// @Param       id path uuid.UUID  true " "
// @Param       request body request.ProductPayload true " "
// @Success     200  {object}  response.WebResponse
// @Failure     400  {object}  response.WebResponse
// @Security		ApiKeyAuth
// @in 					header
// @name 				Authorization
// @Router       /api/admin/products/{id} [put]
func (controller *ProductController) Update(ctx *fiber.Ctx) error {

	currentUserId := helpers.ParseUint(ctx.Locals("userId").(string))

	productId := helpers.ParseUUID(ctx.Params("id"))

	var productModel request.ProductPayload
	err := ctx.BodyParser(&productModel)
	exception.PanicIfNeeded(err)

	result := controller.ProductService.Update(currentUserId, productId, productModel)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary      Delete product
// @Tags         (Admin) Product
// @Accept       json
// @Produce      json
// @Param        id path int  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/products/{id} [delete]
func (controller *ProductController) Delete(ctx *fiber.Ctx) error {

	currentUserId := helpers.ParseUint(ctx.Locals("userId").(string))

	productId := helpers.ParseUUID(ctx.Params("id"))

	msgResult := controller.ProductService.Delete(currentUserId, productId)

	return ctx.JSON(response.NewWebResponse(nil, msgResult))
}

// @Summary		Search products
// @Tags			(Admin) Product
// @Accept		json
// @Produce		json
// @Param			query query string false " "
// @Param			productTypes query string false " "
// @Param			productCategoryId query int false " "
// @Success		200	{object}	[]response.ProductResponse
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/products [get]
func (controller *ProductController) Search(ctx *fiber.Ctx) error {

	currentUserId := helpers.ParseUint(ctx.Locals("userId").(string))

	query := helpers.ParseNullableString(ctx.Query("query"))

	var productTypes *[]string = nil
	productTypeQuery := ctx.Queries()["productTypes"]
	if productTypeQuery != "" {
		productTypesQueries := strings.Split(productTypeQuery, ",")
		productTypes = &productTypesQueries
	}

	productCategoryId := helpers.ParseNullableUint(ctx.Query("productCategoryId"))

	page := helpers.ParseNullableInt(ctx.Query("page"))
	pageSize := helpers.ParseNullableInt(ctx.Query("pageSize"))

	result := controller.ProductService.Search(currentUserId, query, productTypes, productCategoryId, page, pageSize)

	return ctx.JSON(response.NewWebResponse(result))
}

// func (controller *ProductController) BrowseProductSku(ctx *fiber.Ctx) error {

// 	currentUserId := helpers.ParseUint(ctx.Locals("userId").(string))

// 	query := helpers.ParseNullableString(ctx.Query("query"))

// 	var productTypes *[]string = nil
// 	productTypeQuery := ctx.Queries()["productTypes"]
// 	if productTypeQuery != "" {
// 		productTypesQueries := strings.Split(productTypeQuery, ",")
// 		productTypes = &productTypesQueries
// 	}

// 	isInventoryOnly := helpers.ParseNullableBool(ctx.Query("isInventoryOnly"))

// 	// productCategoryId := helpers.ParseNullableUniqueidentifier(ctx.Query("productCategoryId"))

// 	page := helpers.ParseNullableInt(ctx.Query("page"))
// 	pageSize := helpers.ParseNullableInt(ctx.Query("pageSize"))

// 	result := controller.ProductService.SearchProductSku(currentUserId, outletId, query, productTypes, isInventoryOnly, productCategoryId, brandId, page, pageSize)

// 	return ctx.JSON(response.NewWebResponse(result))
// }

// @Summary		Get product by id
// @Tags			(Admin) Product
// @Accept		json
// @Produce		json
// @Param			id path int true " "
// @Success		200	{object}	response.ProductResponse
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/products/{id} [get]
func (controller *ProductController) GetById(ctx *fiber.Ctx) error {

	currentUserId := helpers.ParseUint(ctx.Locals("userId").(string))

	productId := helpers.ParseUUID(ctx.Params("id"))

	result := controller.ProductService.GetById(currentUserId, productId)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary     Create Variant Option
// @Tags        (Admin) Product Variant Option
// @Accept      json
// @Produce     json
// @Param       request body request.VariantOptionPayload true " "
// @Success     200  {object}  response.WebResponse
// @Failure     400  {object}  response.WebResponse
// @Security		ApiKeyAuth
// @in 					header
// @name 				Authorization
// @Router       /api/admin/products/variant-options [post]
func (controller *ProductController) CreateVariantOptions(ctx *fiber.Ctx) error {

	currentUserId := helpers.ParseUint(ctx.Locals("userId").(string))

	var optionModel request.VariantOptionPayload
	err := ctx.BodyParser(&optionModel)
	exception.PanicIfNeeded(err)

	result := controller.ProductService.CreateVariantOptions(currentUserId, optionModel)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary     Get Variant Options
// @Tags        (Admin) Product Variant Option
// @Accept      json
// @Produce     json
// @Success     200  {object}  []response.VariantOption
// @Failure     400  {object}  response.WebResponse
// @Security		ApiKeyAuth
// @in 					header
// @name 				Authorization
// @Router       /api/admin/products/variant-options [get]
func (controller *ProductController) GetVariantOptions(ctx *fiber.Ctx) error {

	currentUserId := helpers.ParseUint(ctx.Locals("userId").(string))

	result := controller.ProductService.GetVariantOptions(currentUserId)

	return ctx.JSON(response.NewWebResponse(result))
}
