package app

import (
	"api.mijkomp.com/controller"

	_ "api.mijkomp.com/docs"
	"api.mijkomp.com/exception"
	"api.mijkomp.com/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

func NewFiberConfig() fiber.Config {
	return fiber.Config{
		BodyLimit:    10 * 1024 * 1024,
		ErrorHandler: exception.ErrorHandler,
	}
}

func CreateServer(
	userService service.UserService,
	productCategoryService service.ProductCategoryService,
	productService service.ProductService,
	db *gorm.DB,
) *fiber.App {

	app := fiber.New(NewFiberConfig())

	var ConfigDefault = recover.Config{
		Next:              nil,
		EnableStackTrace:  true,
		StackTraceHandler: recover.ConfigDefault.StackTraceHandler,
	}
	app.Use(recover.New(ConfigDefault))

	// cors
	app.Use(cors.New())

	// swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	userController := controller.NewUserController(&userService)
	userController.Route(app)
	controller.NewAuthController(&userService).Route(app)
	controller.NewProductCategoryController(&userService, &productCategoryService).Route(app)
	controller.NewProductController(&userService, &productService).Route(app)

	return app
}
