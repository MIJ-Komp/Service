package admin

import (
	"api.mijkomp.com/helpers"
	"api.mijkomp.com/middleware"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/service"
	"github.com/gofiber/fiber/v2"
)

type DashboardController struct {
	UserService      service.UserService
	DashboardService service.DashboardService
}

func NewDashboardController(userService *service.UserService, dashboardService *service.DashboardService) *DashboardController {
	return &DashboardController{UserService: *userService, DashboardService: *dashboardService}
}

func (controller *DashboardController) Route(app *fiber.App) {
	dashboard := app.Group("/api/admin/dashboard", middleware.AuthMiddleware(controller.UserService))
	dashboard.Get("/", controller.GetSummary)
}

// @Summary      Get Summary
// @Description  Get dashboard summary with optional date range filter
// @Tags         Dashboard
// @Accept       json
// @Produce      json
// @Param        fromDate query string false "Start date for filtering (format: 2006-01-02)"
// @Param        toDate query string false "End date for filtering (format: 2006-01-02)"
// @Success      200  {object}  response.WebResponse{data=response.Dashboard}
// @Failure      400  {object}  response.WebResponse
// @Security     ApiKeyAuth
// @Router       /api/admin/dashboard [get]
func (controller *DashboardController) GetSummary(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	fromDate := helpers.ParseDate(ctx.Query("fromDate"))
	toDate := helpers.ParseDate(ctx.Query("toDate"))
	result := controller.DashboardService.GetSummary(currentUserId, fromDate, toDate)

	return ctx.JSON(response.NewWebResponse(result))
}
