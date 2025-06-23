package customer

import (
	"api.mijkomp.com/exception"
	"api.mijkomp.com/middleware"
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/service"
	"github.com/gofiber/fiber/v2"
)

type PaymentController struct {
	PaymentService service.PaymentService
	UserService    service.UserService
}

func NewPaymentController(paymentService *service.PaymentService) *PaymentController {
	return &PaymentController{PaymentService: *paymentService}
}

func (controller *PaymentController) Route(app *fiber.App) {
	payment := app.Group("/api/payment")
	payment.Post("/xendit/updatePayment", middleware.PaymentMiddlewareXendit(), controller.UpdatePayment)
}

func (controller *PaymentController) UpdatePayment(ctx *fiber.Ctx) error {

	var payload request.PaymentNotification
	err := ctx.BodyParser(&payload)
	exception.PanicIfNeeded(err)

	controller.PaymentService.Update(payload)

	return ctx.JSON(response.NewWebResponse(nil, "Ok"))
}
