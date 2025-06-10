//go:build wireinject
// +build wireinject

package main

import (
	"api.mijkomp.com/app"
	"api.mijkomp.com/config"
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

var productSet = wire.NewSet(
	repository_impl.NewProductRepository,
	wire.Bind(new(repository.ProductRepository), new(*repository_impl.ProductRepositoryImpl)),
	service_impl.NewProductService,
	wire.Bind(new(service.ProductService), new(*service_impl.ProductServiceImpl)),
)

func InitializedServer() *fiber.App {
	wire.Build(
		config.NewDB,
		userSet,
		productCategorySet,
		productSet,
		validatorSet,
		app.CreateServer,
	)
	return nil
}
