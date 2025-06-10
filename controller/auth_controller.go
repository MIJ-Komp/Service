package controller

import (
	"api.mijkomp.com/exception"
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/service"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	UserService service.UserService
}

func NewAuthController(userService *service.UserService) *AuthController {

	return &AuthController{UserService: *userService}
}

func (controller *AuthController) Route(app *fiber.App) {
	auth := app.Group("/api/auth")
	auth.Post("/login", controller.Login)
}

// @Summary      Login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body request.LoginUserPayload  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      404  {object}  response.WebResponse
// @Failure      500  {object}  response.WebResponse
// @Router       /api/auth/login [post]
func (controller *AuthController) Login(ctx *fiber.Ctx) error {
	var userModel request.LoginUserPayload
	err := ctx.BodyParser(&userModel)
	exception.PanicIfNeeded(err)

	token := controller.UserService.LoginUser(userModel)

	return ctx.JSON(response.NewWebResponse(token))
}
