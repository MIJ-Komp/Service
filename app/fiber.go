package app

import (
	"api.mijkomp.com/controller/admin"
	"api.mijkomp.com/controller/customer"
	"api.mijkomp.com/exception"
	"api.mijkomp.com/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	"gorm.io/gorm"
)

// pastikan kamu pakai fiber-swagger adapter

// import (
// 	"api.mijkomp.com/controller/admin"
// 	_ "api.mijkomp.com/docs"
// 	"api.mijkomp.com/exception"
// 	"api.mijkomp.com/service"
// 	"github.com/gofiber/adaptor/v2"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/fiber/v2/middleware/cors"
// 	httpSwagger "github.com/swaggo/http-swagger"
// 	"gorm.io/gorm"
// )

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
	componentTypeService service.ComponentTypeService,
	compatibilityRuleService service.CompatibilityRuleService,
	menuService service.MenuService,
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

	app.Get("/swagger/admin/*", swagger.New(swagger.Config{
		URL: "/docs/admin/swagger.json",
	}))

	app.Get("/swagger/customer/*", swagger.New(swagger.Config{
		URL: "/docs/customer/swagger.json",
	}))
	app.Static("/docs", "./docs")

	userController := admin.NewUserController(&userService)
	userController.Route(app)
	admin.NewAuthController(&userService).Route(app)
	admin.NewProductCategoryController(&userService, &productCategoryService).Route(app)
	admin.NewProductController(&userService, &productService).Route(app)
	admin.NewComponentTypeController(&userService, &componentTypeService).Route(app)
	admin.NewCompatibilityRuleController(&userService, &compatibilityRuleService).Route(app)
	admin.NewMenuController(&userService, &menuService).Route(app)

	customer.NewProductCategoryController(&userService, &productCategoryService).Route(app)
	customer.NewProductController(&userService, &productService).Route(app)

	return app
}
