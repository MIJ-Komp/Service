package customer

import (
	"api.mijkomp.com/helpers"
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
	compatibilityRule := app.Group("/api/compatibility-rules")
	compatibilityRule.Get("/", controller.Search)
	compatibilityRule.Get("/:id", controller.GetById)
}

// @Summary		Search compatibility rules
// @Tags			Compatibility Rules
// @Accept		json
// @Produce		json
// @Param			sourceComponentTypeCode query string false " "
// @Param			targetComponentTypeCode query string false " "
// @Success		200	{object}	response.CompatibilityRule
// @Failure		404	{object}	response.WebResponse
// @Router		/api/compatibility-rules [get]
func (controller *CompatibilityRuleController) Search(ctx *fiber.Ctx) error {
	var sourceComponentTypeCode *string = nil
	var targetComponentTypeCode *string = nil

	if str := ctx.Query("sourceComponentTypeCode"); str != "" {
		sourceComponentTypeCode = &str
	}

	if str := ctx.Query("targetComponentTypeCode"); str != "" {
		targetComponentTypeCode = &str
	}

	result := controller.CompatibilityRuleService.Search(0, sourceComponentTypeCode, targetComponentTypeCode)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary		Get compatibility rules by id
// @Tags			Compatibility Rules
// @Accept		json
// @Produce		json
// @Param			id path int true " "
// @Success		200	{object}	response.CompatibilityRule
// @Failure		404	{object}	response.WebResponse
// @Router		/api/compatibility-rules/{id} [get]
func (controller *CompatibilityRuleController) GetById(ctx *fiber.Ctx) error {
	id := helpers.ParseUint(ctx.Params("id"))

	result := controller.CompatibilityRuleService.GetById(0, id)

	return ctx.JSON(response.NewWebResponse(result))
}
