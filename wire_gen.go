// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"api.mijkomp.com/app"
	"api.mijkomp.com/config"
	"api.mijkomp.com/database"
	"api.mijkomp.com/repository"
	"api.mijkomp.com/repository/repository_impl"
	"api.mijkomp.com/service"
	"api.mijkomp.com/service/service_impl"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

// Injectors from injector.go:

func InitializedServer() *fiber.App {
	userRepositoryImpl := repository_impl.NewUserRepository()
	validate := config.NewValidator()
	db := database.NewDB()
	userServiceImpl := service_impl.NewUserService(userRepositoryImpl, validate, db)
	productCategoryRepositoryImpl := repository_impl.NewProductCategoryRepository()
	productCategoryServiceImpl := service_impl.NewProductCategoryService(productCategoryRepositoryImpl, validate, db)
	brandRepositoryImpl := repository_impl.NewBrandRepository()
	brandServiceImpl := service_impl.NewBrandService(brandRepositoryImpl, validate, db)
	productRepositoryImpl := repository_impl.NewProductRepository()
	productServiceImpl := service_impl.NewProductService(productRepositoryImpl, db)
	componentTypeRepositoryImpl := repository_impl.NewComponentTypeRepository()
	componentTypeServiceImpl := service_impl.NewComponentTypeService(componentTypeRepositoryImpl, validate, db)
	compatibilityRuleRepositoryImpl := repository_impl.NewCompatibilityRuleRepository()
	compatibilityRuleServiceImpl := service_impl.NewCompatibilityRuleService(compatibilityRuleRepositoryImpl, validate, db)
	menuRepositoryImpl := repository_impl.NewMenuRepository()
	menuServiceImpl := service_impl.NewMenuService(menuRepositoryImpl, validate, db)
	orderRepositoryImpl := repository_impl.NewOrderRepository()
	paymentRepositoryImpl := repository_impl.NewPaymentRepository()
	paymentServiceImpl := service_impl.NewPaymentService(paymentRepositoryImpl, orderRepositoryImpl, db)
	orderServiceImpl := service_impl.NewOrderService(orderRepositoryImpl, productServiceImpl, productRepositoryImpl, paymentServiceImpl, db)
	dashboardServiceImpl := service_impl.NewDashboardService(db)
	fiberApp := app.CreateServer(userServiceImpl, productCategoryServiceImpl, brandServiceImpl, productServiceImpl, componentTypeServiceImpl, compatibilityRuleServiceImpl, menuServiceImpl, orderServiceImpl, dashboardServiceImpl, db)
	return fiberApp
}

// injector.go:

var validatorSet = wire.NewSet(config.NewValidator)

var userSet = wire.NewSet(repository_impl.NewUserRepository, wire.Bind(new(repository.UserRepository), new(*repository_impl.UserRepositoryImpl)), service_impl.NewUserService, wire.Bind(new(service.UserService), new(*service_impl.UserServiceImpl)))

var productCategorySet = wire.NewSet(repository_impl.NewProductCategoryRepository, wire.Bind(new(repository.ProductCategoryRepository), new(*repository_impl.ProductCategoryRepositoryImpl)), service_impl.NewProductCategoryService, wire.Bind(new(service.ProductCategoryService), new(*service_impl.ProductCategoryServiceImpl)))

var brandSet = wire.NewSet(repository_impl.NewBrandRepository, wire.Bind(new(repository.BrandRepository), new(*repository_impl.BrandRepositoryImpl)), service_impl.NewBrandService, wire.Bind(new(service.BrandService), new(*service_impl.BrandServiceImpl)))

var productSet = wire.NewSet(repository_impl.NewProductRepository, wire.Bind(new(repository.ProductRepository), new(*repository_impl.ProductRepositoryImpl)), service_impl.NewProductService, wire.Bind(new(service.ProductService), new(*service_impl.ProductServiceImpl)))

var componentTypeSet = wire.NewSet(repository_impl.NewComponentTypeRepository, wire.Bind(new(repository.ComponentTypeRepository), new(*repository_impl.ComponentTypeRepositoryImpl)), service_impl.NewComponentTypeService, wire.Bind(new(service.ComponentTypeService), new(*service_impl.ComponentTypeServiceImpl)))

var compatibilityRuleSet = wire.NewSet(repository_impl.NewCompatibilityRuleRepository, wire.Bind(new(repository.CompatibilityRuleRepository), new(*repository_impl.CompatibilityRuleRepositoryImpl)), service_impl.NewCompatibilityRuleService, wire.Bind(new(service.CompatibilityRuleService), new(*service_impl.CompatibilityRuleServiceImpl)))

var menuSet = wire.NewSet(repository_impl.NewMenuRepository, wire.Bind(new(repository.MenuRepository), new(*repository_impl.MenuRepositoryImpl)), service_impl.NewMenuService, wire.Bind(new(service.MenuService), new(*service_impl.MenuServiceImpl)))

var orderSet = wire.NewSet(repository_impl.NewOrderRepository, wire.Bind(new(repository.OrderRepository), new(*repository_impl.OrderRepositoryImpl)), service_impl.NewOrderService, wire.Bind(new(service.OrderService), new(*service_impl.OrderServiceImpl)))

var paymentSet = wire.NewSet(repository_impl.NewPaymentRepository, wire.Bind(new(repository.PaymentRepository), new(*repository_impl.PaymentRepositoryImpl)), service_impl.NewPaymentService, wire.Bind(new(service.PaymentService), new(*service_impl.PaymentServiceImpl)))

var dashboardSet = wire.NewSet(service_impl.NewDashboardService, wire.Bind(new(service.DashboardService), new(*service_impl.DashboardServiceImpl)))
