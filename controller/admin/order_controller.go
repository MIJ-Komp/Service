package admin

import (
	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers"
	"api.mijkomp.com/middleware"
	"api.mijkomp.com/models/enum"
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/service"
	"github.com/gofiber/fiber/v2"
)

type OrderController struct {
	UserService  service.UserService
	OrderService service.OrderService
}

func NewOrderController(userService *service.UserService, orderService *service.OrderService) *OrderController {
	return &OrderController{UserService: *userService, OrderService: *orderService}
}

func (controller *OrderController) Route(app *fiber.App) {
	order := app.Group("/api/admin/orders", middleware.AuthMiddleware(controller.UserService))
	order.Get("/", controller.Search)
	order.Get("/:id", controller.GetById)
	order.Put("/:id/update-status", controller.UpdateStatus)
	order.Put("/:id/update-shipping-info", controller.UpdateShippingInfo)
}

// @Summary		Search order
// @Tags			Order
// @Accept		json
// @Produce		json
// @Param			query query string false " "
// @Param			status query string false " "
// @Param			fromDate query string false " "
// @Param			toDate query string false " "
// @Param			page query int true " "
// @Param			pageSize query int true " "
// @Success		200	{object}	response.Order
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/orders [get]
func (controller *OrderController) Search(ctx *fiber.Ctx) error {

	query := helpers.ParseNullableString(ctx.Query("query"))

	statusString := ctx.Query("status")
	var status *enum.EOrderStatus
	if statusString != "" {
		statusEnum := enum.EOrderStatus(statusString)
		if !statusEnum.IsValid() {
			panic(exception.NewValidationError("Status invalid."))
		}
		status = &statusEnum
	}

	fromDate := helpers.ParseNullableTime(ctx.Query("fromDate"))
	toDate := helpers.ParseNullableTime(ctx.Query("toDate"))
	page := helpers.ParseInt(ctx.Query("page"))
	pageSize := helpers.ParseInt(ctx.Query("pageSize"))

	result := controller.OrderService.Search(query, status, fromDate, toDate, page, pageSize)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary		Get order by id
// @Tags			Order
// @Accept		json
// @Produce		json
// @Param			id path string true " "
// @Success		200	{object}	response.Order
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/orders/{id} [get]
func (controller *OrderController) GetById(ctx *fiber.Ctx) error {

	id := helpers.ParseUUID(ctx.Params("id"))

	result := controller.OrderService.GetById(&id, nil)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary		Update status order
// @Tags			Order
// @Accept		json
// @Produce		json
// @Param			id path string true " "
// @Param     request body request.UpdateOrderStatusByAdmin  true " "
// @Success		200	{object}	response.Order
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/orders/{id}/update-status [put]
func (controller *OrderController) UpdateStatus(ctx *fiber.Ctx) error {

	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))
	id := helpers.ParseUUID(ctx.Params("id"))

	var payload request.UpdateOrderStatusByAdmin
	err := ctx.BodyParser(&payload)
	exception.PanicIfNeeded(err)

	result := controller.OrderService.UpdateStatus(currentUserId, id, payload)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary		Update shipping info order
// @Tags			Order
// @Accept		json
// @Produce		json
// @Param			id path string true " "
// @Param     request body request.UpdateOrderShippingByAdmin  true " "
// @Success		200	{object}	response.Order
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/admin/orders/{id}/update-shipping-info [put]
func (controller *OrderController) UpdateShippingInfo(ctx *fiber.Ctx) error {

	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))
	id := helpers.ParseUUID(ctx.Params("id"))

	var payload request.UpdateOrderShippingByAdmin
	err := ctx.BodyParser(&payload)
	exception.PanicIfNeeded(err)

	result := controller.OrderService.UpdateShippingInfo(currentUserId, id, payload)

	return ctx.JSON(response.NewWebResponse(result))
}
