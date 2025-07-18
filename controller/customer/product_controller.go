package customer

import (
	"strings"

	"api.mijkomp.com/helpers"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
// @Param  		ids query []string false "Array of Ids (1,2,3)"
// @Param			productTypes query string false " "
// @Param  		componentTypeIds query string false "Array of IDs (1,2,3)"
// @Param  		productCategoryIds query string false "Array of IDs (1,2,3)"
// @Param  		brandIds query string false "Array of IDs (1,2,3)"
// @Param  		tag query string false "String"
// @Param  		minPrice query float64 false "money"
// @Param  		maxPrice query float64 false "money"
// @Param  		isInStockOnly query bool false "boolean"
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

	var ids *[]uuid.UUID = nil
	productIdQuery := ctx.Queries()["ids"]
	if productIdQuery != "" {
		productIdQueries := strings.Split(productIdQuery, ",")

		uuidIds := []uuid.UUID{}
		for _, strId := range productIdQueries {
			uuidIds = append(uuidIds, helpers.ParseUUID(strId))
		}
		if len(uuidIds) > 0 {
			ids = &uuidIds
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

	var brandIds *[]uint = nil
	brandIdsQuery := ctx.Queries()["brandIds"]
	if brandIdsQuery != "" {
		brandIdsQueries := strings.Split(brandIdsQuery, ",")
		var ids []uint
		for _, brandQuery := range brandIdsQueries {
			id := helpers.ParseUint(brandQuery)
			ids = append(ids, id)
		}
		if len(ids) > 0 {
			brandIds = &ids
		}
	}

	tag := helpers.ParseNullableString(ctx.Query("tag"))

	minPrice := helpers.ParseNullableFloat64(ctx.Query("minPrice"))
	maxPrice := helpers.ParseNullableFloat64(ctx.Query("maxPrice"))

	isInStockOnly := helpers.ParseNullableBool(ctx.Query("isInStockOnly"))

	isActive := true
	isShowOnlyInMarketPlace := false

	page := helpers.ParseNullableInt(ctx.Query("page"))
	pageSize := helpers.ParseNullableInt(ctx.Query("pageSize"))

	result := controller.ProductService.Search(0, query, ids, productTypes, componentTypeIds, productCategoryIds, brandIds, tag, minPrice, maxPrice, isInStockOnly, &isActive, &isShowOnlyInMarketPlace, page, pageSize)

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
