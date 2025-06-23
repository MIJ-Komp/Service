package customer

import (
	"api.mijkomp.com/exception"
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
	order := app.Group("/api/orders")
	order.Post("/", controller.Create)
	// order.Put("/:id", controller.Update)
	// order.Delete("/:id", controller.Delete)
	// order.Get("/", controller.Search)
	order.Get("/:code", controller.GetById)
}

// @Summary      Create order
// @Tags         Order
// @Accept       json
// @Produce      json
// @Param        request body request.Order  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router       /api/orders [post]
func (controller *OrderController) Create(ctx *fiber.Ctx) error {

	var payload request.Order
	err := ctx.BodyParser(&payload)
	exception.PanicIfNeeded(err)

	result := controller.OrderService.Create(payload)

	return ctx.JSON(response.NewWebResponse(result))
}

// // @Summary      Update order
// // @Tags         Order
// // @Accept       json
// // @Produce      json
// // @Param        id path int  true " "
// // @Param        request body request.Order  true " "
// // @Success      200  {object}  response.WebResponse
// // @Failure      400  {object}  response.WebResponse
// // @Security	ApiKeyAuth
// // @in header
// // @name Authorization
// // @Router       /api/orders/{id} [put]
// func (controller *OrderController) Update(ctx *fiber.Ctx) error {

// 	orderId := helpers.ParseUUID(ctx.Params("id"))

// 	var payload request.Order
// 	err := ctx.BodyParser(&payload)
// 	exception.PanicIfNeeded(err)

// 	result := controller.OrderService.Update(orderId, payload)

// 	return ctx.JSON(response.NewWebResponse(result))
// }

// // @Summary      Delete order
// // @Tags         Order
// // @Accept       json
// // @Produce      json
// // @Param        id path int  true " "
// // @Success      200  {object}  response.WebResponse
// // @Failure      400  {object}  response.WebResponse
// // @Security	ApiKeyAuth
// // @in header
// // @name Authorization
// // @Router       /api/orders/{id} [delete]
// func (controller *OrderController) Delete(ctx *fiber.Ctx) error {

// 	orderId := helpers.ParseUUID(ctx.Params("id"))
// 	result := controller.OrderService.Delete(orderId)

// 	return ctx.JSON(response.NewWebResponse(nil, result))
// }

// // @Summary		Search order
// // @Tags			Order
// // @Accept		json
// // @Produce		json
// // @Param			query query string false " "
// // @Param			status query string false " "
// // @Param			fromDate query string false " "
// // @Param			toDate query time.Time false " "
// // @Param			page query int true " "
// // @Param			pageSize query int true " "
// // @Success		200	{object}	response.Order
// // @Failure		404	{object}	response.WebResponse
// // @Security	ApiKeyAuth
// // @in header
// // @name Authorization
// // @Router		/api/orders [get]
// func (controller *OrderController) Search(ctx *fiber.Ctx) error {

// 	query := helpers.ParseNullableString(ctx.Query("query"))

// 	statusString := ctx.Query("status")
// 	var status *enum.EOrderStatus
// 	if statusString != "" {
// 		statusEnum := enum.EOrderStatus(statusString)
// 		if !statusEnum.IsValid() {
// 			panic(exception.NewValidationError("Status invalid."))
// 		}
// 		status = &statusEnum
// 	}

// 	fromDate := helpers.ParseNullableTime(ctx.Query("fromDate"))
// 	toDate := helpers.ParseNullableTime(ctx.Query("toDate"))
// 	page := helpers.ParseInt(ctx.Query("page"))
// 	pageSize := helpers.ParseInt(ctx.Query("pageSize"))

// 	result := controller.OrderService.Search(query, status, fromDate, toDate, page, pageSize)

// 	return ctx.JSON(response.NewWebResponse(result))
// }

// @Summary		Get order by code
// @Tags			Order
// @Accept		json
// @Produce		json
// @Param			code path string true " "
// @Success		200	{object}	response.Order
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/orders/{code} [get]
func (controller *OrderController) GetById(ctx *fiber.Ctx) error {

	code := ctx.Params("code")

	result := controller.OrderService.GetById(nil, &code)

	return ctx.JSON(response.NewWebResponse(result))
}
