package admin

import (
	"fmt"
	"strings"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers"
	"api.mijkomp.com/helpers/logger"
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
	product.Put("/change-component", middleware.AuthMiddleware(controller.UserService), controller.UpdateComponent)
	product.Put("/:id", middleware.AuthMiddleware(controller.UserService), controller.Update)
	product.Delete("/:id", middleware.AuthMiddleware(controller.UserService), controller.Delete)
	product.Get("/", middleware.AuthMiddleware(controller.UserService), controller.Search)
	product.Get("/browse-product-sku", middleware.AuthMiddleware(controller.UserService), controller.BrowseProductSku)
	product.Get("/:id", middleware.AuthMiddleware(controller.UserService), controller.GetById)

	variantOption := app.Group("/api/admin/variant-options")
	variantOption.Post("/", middleware.AuthMiddleware(controller.UserService), controller.CreateVariantOptions)
	variantOption.Get("/", middleware.AuthMiddleware(controller.UserService), controller.GetVariantOptions)
}

// @Summary     Create Product
// @Tags        Product
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

	logger.LogInfo(fmt.Sprintf("[Admin %d] Mencoba membuat produk baru", currentUserId))

	var productModel request.ProductPayload
	err := ctx.BodyParser(&productModel)
	if err != nil {
		panic(exception.NewValidationError(err.Error()))
	}

	result := controller.ProductService.Create(currentUserId, productId, productModel)

	logger.LogInfo(fmt.Sprintf("[Admin %d] Berhasil membuat produk baru", currentUserId))

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary     Update Product
// @Tags        Product
// @Accept      json
// @Produce     json
// @Param       id path string  true " "
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

	logger.LogInfo(fmt.Sprintf("[Admin %d] Mencoba mengupdate produk dengan ID %s", currentUserId, productId))

	var productModel request.ProductPayload
	err := ctx.BodyParser(&productModel)
	exception.PanicIfNeeded(err)

	result := controller.ProductService.Update(currentUserId, productId, productModel)

	logger.LogInfo(fmt.Sprintf("[Admin %d] Berhasil mengupdate produk dengan ID %s", currentUserId, productId))

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary     change component (for product bundle)
// @Tags        Product
// @Accept      json
// @Produce     json
// @Param       request body request.ChangeComponent true " "
// @Success     200  {object}  response.WebResponse
// @Failure     400  {object}  response.WebResponse
// @Security		ApiKeyAuth
// @in 					header
// @name 				Authorization
// @Router       /api/admin/products/change-component [put]
func (controller *ProductController) UpdateComponent(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUint(ctx.Locals("userId").(string))

	var payload request.ChangeComponent
	err := ctx.BodyParser(&payload)
	exception.PanicIfNeeded(err)

	logger.LogInfo(fmt.Sprintf("[Admin %d] mengubah komponent pada semua produk bundle", currentUserId))

	result := controller.ProductService.ChangeComponent(currentUserId, payload)

	logger.LogInfo(fmt.Sprintf("[Admin %d] berhasil mengubah komponent pada semua produk bundle", currentUserId))

	return ctx.JSON(response.NewWebResponse(nil, result))
}

// @Summary      Delete product
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id path string  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/products/{id} [delete]
func (controller *ProductController) Delete(ctx *fiber.Ctx) error {

	currentUserId := helpers.ParseUint(ctx.Locals("userId").(string))

	productId := helpers.ParseUUID(ctx.Params("id"))

	logger.LogInfo(fmt.Sprintf("[Admin %d] Mencoba menghapus produk dengan ID %s", currentUserId, productId))

	msgResult := controller.ProductService.Delete(currentUserId, productId)

	logger.LogInfo(fmt.Sprintf("[Admin %d] Berhasil menghapus produk dengan ID %s", currentUserId, productId))

	return ctx.JSON(response.NewWebResponse(nil, msgResult))
}

// @Summary		Search products
// @Tags			Product
// @Accept		json
// @Produce		json
// @Param			query query string false " "
// @Param			productTypes query string false " "
// @Param  		productCategoryIds query []int false "Array of IDs" collectionFormat(csv)
// @Param  		componentTypeIds query []int false "Array of IDs" collectionFormat(csv)
// @Param			isActive query bool false " "
// @Param			isShowOnlyInMarketplace query bool false " "
// @Param			page query int false " "
// @Param			pageSize query int false " "
// @Success		200	{object}	[]response.ProductResponse
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/products [get]
func (controller *ProductController) Search(ctx *fiber.Ctx) error {

	currentUserId := helpers.ParseUint(ctx.Locals("userId").(string))

	logger.LogInfo(fmt.Sprintf("[Admin %d] Mencari produk", currentUserId))

	query := helpers.ParseNullableString(ctx.Query("query"))

	var productTypes *[]string = nil
	productTypeQuery := ctx.Queries()["productTypes"]
	if productTypeQuery != "" {
		productTypesQueries := strings.Split(productTypeQuery, ",")
		productTypes = &productTypesQueries
	}

	// productCategoryId := helpers.ParseNullableUint(ctx.Query("productCategoryId"))
	var productCategoryIds *[]uint = nil
	productCategoryIdsQuery := ctx.Queries()["productCategoryIds"]

	if productCategoryIdsQuery != "" {
		productCategoryQueries := strings.Split(productCategoryIdsQuery, ",")
		var ids []uint
		for _, productCategoryQuery := range productCategoryQueries {
			id := helpers.ParseUint(productCategoryQuery)
			ids = append(ids, id)
		}
		if len(ids) > 0 {
			productCategoryIds = &ids
		}
	}

	var componentTypeIds *[]uint = nil
	if cIdsQuery := ctx.Queries()["componentTypeIds"]; cIdsQuery != "" {
		cIds := strings.Split(cIdsQuery, ",")
		ids := []uint{}

		for _, cid := range cIds {
			ids = append(ids, helpers.ParseUint(cid))
		}
		if len(ids) > 0 {
			componentTypeIds = &ids
		}
	}

	isActive := helpers.ParseNullableBool(ctx.Query("isActive"))
	isShowOnlyInMarketPlace := helpers.ParseNullableBool(ctx.Query("isShowOnlyInMarketPlace"))

	page := helpers.ParseNullableInt(ctx.Query("page"))
	pageSize := helpers.ParseNullableInt(ctx.Query("pageSize"))

	result := controller.ProductService.Search(currentUserId, query, nil, productTypes, componentTypeIds, productCategoryIds, nil, nil, nil, nil, nil, isActive, isShowOnlyInMarketPlace, page, pageSize)

	logger.LogInfo(fmt.Sprintf("[Admin %d] Berhasil mendapatkan daftar produk", currentUserId))

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary		Browse products sku
// @Tags			Product
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
// @Router		/api/admin/products/browse-product-sku [get]
func (controller *ProductController) BrowseProductSku(ctx *fiber.Ctx) error {

	currentUserId := helpers.ParseUint(ctx.Locals("userId").(string))

	logger.LogInfo(fmt.Sprintf("[Admin %d] Mencari SKU produk", currentUserId))

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

	result := controller.ProductService.BrowseProductSku(currentUserId, query, productTypes, productCategoryId, page, pageSize)

	logger.LogInfo(fmt.Sprintf("[Admin %d] Berhasil mendapatkan daftar SKU produk", currentUserId))

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary		Get product by id
// @Tags			Product
// @Accept		json
// @Produce		json
// @Param			id path string true " "
// @Success		200	{object}	response.ProductResponse
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/products/{id} [get]
func (controller *ProductController) GetById(ctx *fiber.Ctx) error {

	currentUserId := helpers.ParseUint(ctx.Locals("userId").(string))

	productId := helpers.ParseUUID(ctx.Params("id"))

	logger.LogInfo(fmt.Sprintf("[Admin %d] Mencoba mendapatkan produk dengan ID %s", currentUserId, productId))

	result := controller.ProductService.GetById(currentUserId, productId)

	logger.LogInfo(fmt.Sprintf("[Admin %d] Berhasil mendapatkan produk dengan ID %s", currentUserId, productId))

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary     Create Variant Option
// @Tags        Product Variant Option
// @Accept      json
// @Produce     json
// @Param       request body request.VariantOptionPayload true " "
// @Success     200  {object}  response.WebResponse
// @Failure     400  {object}  response.WebResponse
// @Security		ApiKeyAuth
// @in 					header
// @name 				Authorization
// @Router       /api/admin/variant-options [post]
func (controller *ProductController) CreateVariantOptions(ctx *fiber.Ctx) error {

	currentUserId := helpers.ParseUint(ctx.Locals("userId").(string))

	logger.LogInfo(fmt.Sprintf("[Admin %d] Mencoba membuat opsi varian baru", currentUserId))

	var optionModel request.VariantOptionPayload
	err := ctx.BodyParser(&optionModel)
	exception.PanicIfNeeded(err)

	result := controller.ProductService.CreateVariantOptions(currentUserId, optionModel)

	logger.LogInfo(fmt.Sprintf("[Admin %d] Berhasil membuat opsi varian baru", currentUserId))

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary     Get Variant Options
// @Tags        Product Variant Option
// @Accept      json
// @Produce     json
// @Success     200  {object}  []response.VariantOption
// @Failure     400  {object}  response.WebResponse
// @Security		ApiKeyAuth
// @in 					header
// @name 				Authorization
// @Router       /api/admin/variant-options [get]
func (controller *ProductController) GetVariantOptions(ctx *fiber.Ctx) error {

	currentUserId := helpers.ParseUint(ctx.Locals("userId").(string))

	logger.LogInfo(fmt.Sprintf("[Admin %d] Mencoba mendapatkan daftar opsi varian", currentUserId))

	result := controller.ProductService.GetVariantOptions(currentUserId)

	logger.LogInfo(fmt.Sprintf("[Admin %d] Berhasil mendapatkan daftar opsi varian", currentUserId))

	return ctx.JSON(response.NewWebResponse(result))
}
