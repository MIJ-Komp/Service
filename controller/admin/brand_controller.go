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

type BrandController struct {
	UserService  service.UserService
	BrandService service.BrandService
}

func NewBrandController(userService *service.UserService, brandService *service.BrandService) *BrandController {
	return &BrandController{UserService: *userService, BrandService: *brandService}
}

func (controller *BrandController) Route(app *fiber.App) {
	brand := app.Group("/api/admin/brands", middleware.AuthMiddleware(controller.UserService))
	brand.Post("/", controller.Create)
	brand.Put("/:id", controller.Update)
	brand.Delete("/:id", controller.Delete)
	brand.Get("/", controller.Search)
	brand.Get("/:id", controller.GetById)
}

// @Summary      Create brand
// @Tags         Brand
// @Accept       json
// @Produce      json
// @Param        request body request.Brand  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/brands [post]
func (controller *BrandController) Create(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	var payload request.Brand
	err := ctx.BodyParser(&payload)
	exception.PanicIfNeeded(err)

	result := controller.BrandService.Create(currentUserId, payload)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary      Update brand
// @Tags         Brand
// @Accept       json
// @Produce      json
// @Param        id path int  true " "
// @Param        request body request.Brand  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/brands/{id} [put]
func (controller *BrandController) Update(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))
	brandId := helpers.ParseUserId(ctx.Params("id"))

	var payload request.Brand
	err := ctx.BodyParser(&payload)
	exception.PanicIfNeeded(err)

	result := controller.BrandService.Update(currentUserId, brandId, payload)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary      Delete brand
// @Tags         Brand
// @Accept       json
// @Produce      json
// @Param        id path int  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/brands/{id} [delete]
func (controller *BrandController) Delete(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))
	brandId := helpers.ParseUserId(ctx.Params("id"))

	result := controller.BrandService.Delete(currentUserId, brandId)

	return ctx.JSON(response.NewWebResponse(nil, result))
}

// @Summary		Search brand
// @Tags			Brand
// @Accept		json
// @Produce		json
// @Param			query query string false " "
// @Success		200	{object}	response.Brand
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/brands [get]
func (controller *BrandController) Search(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	query := ctx.Query("query")
	result := controller.BrandService.Search(currentUserId, &query)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary		Get brand by id
// @Tags			Brand
// @Accept		json
// @Produce		json
// @Param			id path int true " "
// @Success		200	{object}	response.Brand
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/brands/{id} [get]
func (controller *BrandController) GetById(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	id := helpers.ParseUint(ctx.Params("id"))

	result := controller.BrandService.GetById(currentUserId, id)

	return ctx.JSON(response.NewWebResponse(result))
}
