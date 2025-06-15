package customer

import (
	"strings"

	"api.mijkomp.com/helpers"
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

	product := app.Group("/api/products")
	product.Get("/", controller.Search)
	product.Get("/:id", controller.GetById)
}

// @Summary		Search products
// @Tags			Product
// @Accept		json
// @Produce		json
// @Param			query query string false " "
// @Param			productTypes query string false " "
// @Param			productCategoryId query int false " "
// @Param			page query int false " "
// @Param			pageSize query int false " "
// @Success		200	{object}	[]response.ProductResponse
// @Failure		404	{object}	response.WebResponse
// @Router		/api/products [get]
func (controller *ProductController) Search(ctx *fiber.Ctx) error {

	query := helpers.ParseNullableString(ctx.Query("query"))

	var productTypes *[]string = nil
	productTypeQuery := ctx.Queries()["productTypes"]
	if productTypeQuery != "" {
		productTypesQueries := strings.Split(productTypeQuery, ",")
		productTypes = &productTypesQueries
	}

	productCategoryId := helpers.ParseNullableUint(ctx.Query("productCategoryId"))
	isActive := true
	isShowOnlyInMarketPlace := false

	page := helpers.ParseNullableInt(ctx.Query("page"))
	pageSize := helpers.ParseNullableInt(ctx.Query("pageSize"))

	result := controller.ProductService.Search(0, query, productTypes, productCategoryId, &isActive, &isShowOnlyInMarketPlace, page, pageSize)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary		Get product by id
// @Tags			Product
// @Accept		json
// @Produce		json
// @Param			id path string true " "
// @Success		200	{object}	response.ProductResponse
// @Failure		404	{object}	response.WebResponse
// @Router		/api/products/{id} [get]
func (controller *ProductController) GetById(ctx *fiber.Ctx) error {
	productId := helpers.ParseUUID(ctx.Params("id"))

	result := controller.ProductService.GetById(0, productId)

	return ctx.JSON(response.NewWebResponse(result))
}
