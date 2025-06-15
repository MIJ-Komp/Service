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
	// app.Get("/swagger/admin/*", swagger.HandlerDefault)
	// app.Get("/swagger/customer/*", swagger.HandlerDefault)

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

	customer.NewProductCategoryController(&userService, &productCategoryService).Route(app)
	customer.NewProductController(&userService, &productService).Route(app)

	return app
}

// func swaggerHandler(urlPrefix string) fiber.Handler {
// 	return func(c *fiber.Ctx) error {

// 		path := strings.TrimPrefix(c.OriginalURL(), urlPrefix)
// 		url := "http://localhost:5000" + urlPrefix + "/swagger.json" // sesuaikan port dan base path
// 		if path == "/" || path == "" {
// 			return c.Redirect(urlPrefix+"/index.html?url="+url, http.StatusMovedPermanently)
// 		}
// 		return httpSwagger.Handler(
// 			httpSwagger.URL(url),
// 		)(c.Context())
// 	}
// }

// func swaggerHandler(urlPrefix string) fiber.Handler {
// 	return adaptor.HTTPHandler(httpSwagger.Handler(
// 		httpSwagger.URL("http://localhost:5000/swagger/" + urlPrefix + "/swagger.json"),
// 	))
// }

// func swaggerHandler(jsonPath string) fiber.Handler {
// 	return httpSwagger.Handler(
// 		httpSwagger.URL(jsonPath),
// 	)
// }
