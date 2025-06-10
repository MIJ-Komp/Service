package controller

import (
	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers"
	"api.mijkomp.com/middleware"
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/service"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{UserService: *userService}
}

func (controller *UserController) Route(app *fiber.App) {
	user := app.Group("/api/user")
	user.Post("/register", controller.Register)
	user.Get("/me", middleware.AuthMiddleware(controller.UserService), controller.GetMyProfile)
	user.Get("/getByEmail", controller.GetByEmail)
}

// @Summary      Register user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request body request.RegisterUserPayload  true " "
// @Success      200  {object}  response.WebResponse
// @Failure      400  {object}  response.WebResponse
// @Router       /api/user/register [post]
func (controller *UserController) Register(ctx *fiber.Ctx) error {
	var userModel request.RegisterUserPayload
	err := ctx.BodyParser(&userModel)
	exception.PanicIfNeeded(err)

	result := controller.UserService.Create(userModel)

	userLoginModel := request.LoginUserPayload{
		Email:    result.Email,
		Password: userModel.Password,
	}
	token := controller.UserService.LoginUser(userLoginModel)

	return ctx.JSON(response.NewWebResponse(token))
}

// @Summary		Get my profile
// @Tags			User
// @Accept		json
// @Produce		json
// @Success		200	{object}	response.User
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/user/me [get]
func (controller *UserController) GetMyProfile(ctx *fiber.Ctx) error {
	currentUserId := helpers.ParseUserId(ctx.Locals("userId"))

	result := controller.UserService.GetById(currentUserId)

	return ctx.JSON(response.NewWebResponse(result))
}

// @Summary		Get user by email
// @Tags			User
// @Accept		json
// @Produce		json
// @Param			email query string true " "
// @Success		200	{object}	response.User
// @Failure		404	{object}	response.WebResponse
// @Security	ApiKeyAuth
// @in header
// @name Authorization
// @Router		/api/user/getByEmail [get]
func (controller *UserController) GetByEmail(ctx *fiber.Ctx) error {
	email := ctx.Query("email")
	// helpers.PanicIfError(err)

	result := controller.UserService.GetByEmail(email)

	return ctx.JSON(response.NewWebResponse(result))
}
