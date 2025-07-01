//go:build wireinject
// +build wireinject

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

var validatorSet = wire.NewSet(config.NewValidator)

var userSet = wire.NewSet(
	repository_impl.NewUserRepository,
	wire.Bind(new(repository.UserRepository), new(*repository_impl.UserRepositoryImpl)),
	service_impl.NewUserService,
	wire.Bind(new(service.UserService), new(*service_impl.UserServiceImpl)),
)

var productCategorySet = wire.NewSet(
	repository_impl.NewProductCategoryRepository,
	wire.Bind(new(repository.ProductCategoryRepository), new(*repository_impl.ProductCategoryRepositoryImpl)),
	service_impl.NewProductCategoryService,
	wire.Bind(new(service.ProductCategoryService), new(*service_impl.ProductCategoryServiceImpl)),
)

var brandSet = wire.NewSet(
	repository_impl.NewBrandRepository,
	wire.Bind(new(repository.BrandRepository), new(*repository_impl.BrandRepositoryImpl)),
	service_impl.NewBrandService,
	wire.Bind(new(service.BrandService), new(*service_impl.BrandServiceImpl)),
)

var productSet = wire.NewSet(
	repository_impl.NewProductRepository,
	wire.Bind(new(repository.ProductRepository), new(*repository_impl.ProductRepositoryImpl)),
	service_impl.NewProductService,
	wire.Bind(new(service.ProductService), new(*service_impl.ProductServiceImpl)),
)

var componentTypeSet = wire.NewSet(
	repository_impl.NewComponentTypeRepository,
	wire.Bind(new(repository.ComponentTypeRepository), new(*repository_impl.ComponentTypeRepositoryImpl)),
	service_impl.NewComponentTypeService,
	wire.Bind(new(service.ComponentTypeService), new(*service_impl.ComponentTypeServiceImpl)),
)

var compatibilityRuleSet = wire.NewSet(
	repository_impl.NewCompatibilityRuleRepository,
	wire.Bind(new(repository.CompatibilityRuleRepository), new(*repository_impl.CompatibilityRuleRepositoryImpl)),
	service_impl.NewCompatibilityRuleService,
	wire.Bind(new(service.CompatibilityRuleService), new(*service_impl.CompatibilityRuleServiceImpl)),
)

var menuSet = wire.NewSet(
	repository_impl.NewMenuRepository,
	wire.Bind(new(repository.MenuRepository), new(*repository_impl.MenuRepositoryImpl)),
	service_impl.NewMenuService,
	wire.Bind(new(service.MenuService), new(*service_impl.MenuServiceImpl)),
)

var orderSet = wire.NewSet(
	repository_impl.NewOrderRepository,
	wire.Bind(new(repository.OrderRepository), new(*repository_impl.OrderRepositoryImpl)),
	service_impl.NewOrderService,
	wire.Bind(new(service.OrderService), new(*service_impl.OrderServiceImpl)),
)

var paymentSet = wire.NewSet(
	repository_impl.NewPaymentRepository,
	wire.Bind(new(repository.PaymentRepository), new(*repository_impl.PaymentRepositoryImpl)),
	service_impl.NewPaymentService,
	wire.Bind(new(service.PaymentService), new(*service_impl.PaymentServiceImpl)),
)

func InitializedServer() *fiber.App {
	wire.Build(
		database.NewDB,
		userSet,
		productCategorySet,
		brandSet,
		productSet,
		componentTypeSet,
		compatibilityRuleSet,
		menuSet,
		orderSet,
		paymentSet,
		validatorSet,
		app.CreateServer,
	)
	return nil
}
