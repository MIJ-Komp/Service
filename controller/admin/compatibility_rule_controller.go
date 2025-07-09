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

type CompatibilityRuleController struct {
	UserService              service.UserService
	CompatibilityRuleService service.CompatibilityRuleService
}

func NewCompatibilityRuleController(userService *service.UserService, compatibilityRuleService *service.CompatibilityRuleService) *CompatibilityRuleController {
	return &CompatibilityRuleController{UserService: *userService, CompatibilityRuleService: *compatibilityRuleService}
}

func (controller *CompatibilityRuleController) Route(app *fiber.App) {
	compatibilityRule := app.Group("/api/admin/compatibility-rules", middleware.AuthMiddleware(controller.UserService))
	compatibilityRule.Post("/", controller.Create)
	compatibilityRule.Put("/:id", controller.Update)
	compatibilityRule.Delete("/:id", controller.Delete)
	compatibilityRule.Get("/", controller.Search)
	compatibilityRule.Get("/:id", controller.GetById)
}

// @Summary      Create compatibility rules
// @Tags         Compatibility Rules
// @Accept       json
// @Produce      json
// @Param        request body request.CompatibilityRule  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security		 ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/compatibility-rules [post]
func (controller *CompatibilityRuleController) Create(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	var payload request.CompatibilityRule
	err := ctx.BodyParser(&payload)
	exception.PanicIfNeeded(err)

	result := controller.CompatibilityRuleService.Create(currentUserId, payload)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary      Update compatibility rules
// @Tags         Compatibility Rules
// @Accept       json
// @Produce      json
// @Param        id path int  true " "
// @Param        request body request.CompatibilityRule  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/compatibility-rules/{id} [put]
func (controller *CompatibilityRuleController) Update(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))
	compatibilityRuleId := helpers.ParseUserId(ctx.Params("id"))

	var payload request.CompatibilityRule
	err := ctx.BodyParser(&payload)
	exception.PanicIfNeeded(err)

	result := controller.CompatibilityRuleService.Update(currentUserId, compatibilityRuleId, payload)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary      Delete compatibility rules
// @Tags         Compatibility Rules
// @Accept       json
// @Produce      json
// @Param        id path int  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security		 ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/admin/compatibility-rules/{id} [delete]
func (controller *CompatibilityRuleController) Delete(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))
	compatibilityRuleId := helpers.ParseUserId(ctx.Params("id"))

	result := controller.CompatibilityRuleService.Delete(currentUserId, compatibilityRuleId)

	return ctx.JSON(response.NewWebResponse(nil, result))
}

// @Summary		Search compatibility rules
// @Tags			Compatibility Rules
// @Accept		json
// @Produce		json
// @Param			query query string false " "
// @Success		200	{object}	response.CompatibilityRule
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/compatibility-rules [get]
func (controller *CompatibilityRuleController) Search(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	var sourceComponentTypeCode *string = nil
	var targetComponentTypeCode *string = nil

	if str := ctx.Query("sourceComponentTypeCode"); str != "" {
		sourceComponentTypeCode = &str
	}

	if str := ctx.Query("targetComponentTypeCode"); str != "" {
		targetComponentTypeCode = &str
	}

	result := controller.CompatibilityRuleService.Search(currentUserId, sourceComponentTypeCode, targetComponentTypeCode)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary		Get compatibility rules by id
// @Tags			Compatibility Rules
// @Accept		json
// @Produce		json
// @Param			id path int true " "
// @Success		200	{object}	response.CompatibilityRule
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/compatibility-rules/{id} [get]
func (controller *CompatibilityRuleController) GetById(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	id := helpers.ParseUint(ctx.Params("id"))

	result := controller.CompatibilityRuleService.GetById(currentUserId, id)

	return ctx.JSON(response.NewWebResponse(result))
}
