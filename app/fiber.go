package app

import (
	"api.mijkomp.com/controller/admin"
	"api.mijkomp.com/controller/customer"
	"api.mijkomp.com/exception"
	"api.mijkomp.com/middleware"
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
	brandService service.BrandService,
	productService service.ProductService,
	componentTypeService service.ComponentTypeService,
	compatibilityRuleService service.CompatibilityRuleService,
	menuService service.MenuService,
	orderService service.OrderService,
	dashboardService service.DashboardService,
	db *gorm.DB,
) *fiber.App {

	app := fiber.New(NewFiberConfig())

	var ConfigDefault = recover.Config{
		Next:              nil,
		EnableStackTrace:  true,
		StackTraceHandler: recover.ConfigDefault.StackTraceHandler,
	}
	app.Use(recover.New(ConfigDefault))

	// Request ID middleware (must be before logger)
	app.Use(middleware.RequestIDMiddleware())

	// Logger middleware
	app.Use(middleware.LoggerMiddleware())

	// app.Static("/uploads", "./uploads")

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
	admin.NewBrandController(&userService, &brandService).Route(app)
	admin.NewProductCategoryController(&userService, &productCategoryService).Route(app)
	admin.NewProductController(&userService, &productService).Route(app)
	admin.NewComponentTypeController(&userService, &componentTypeService).Route(app)
	admin.NewCompatibilityRuleController(&userService, &compatibilityRuleService).Route(app)
	admin.NewMenuController(&userService, &menuService).Route(app)
	admin.NewOrderController(&userService, &orderService).Route(app)
	admin.NewFileController(&userService).Route(app)
	admin.NewDashboardController(&userService, &dashboardService).Route(app)

	customer.NewProductCategoryController(&userService, &productCategoryService).Route(app)
	customer.NewProductController(&userService, &productService).Route(app)
	customer.NewMenuController(&userService, &menuService).Route(app)
	customer.NewCompatibilityRuleController(&userService, &compatibilityRuleService).Route(app)
	customer.NewOrderController(&userService, &orderService).Route(app)
	customer.NewFileController(&userService).Route(app)
	return app
}
