package config

import (
	"os"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers"
	"api.mijkomp.com/models/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	exception.PanicIfNeeded(err)

	// Auto migrate entity
	db.AutoMigrate(
		&entity.User{},
		&entity.UserVerification{},
		&entity.UserToken{},
		&entity.ProductCategory{},
		&entity.Product{},
		&entity.ProductSku{},
		// &entity.ProductSkuDetail{},
		&entity.ProductGroupItem{},
		&entity.ProductVariantOption{},
		&entity.ProductVariantOptionValue{},
		&entity.ProductSkuVariant{},
		&entity.VariantOption{},
	)

	seedData(db)
	return db
}

func seedData(db *gorm.DB) {

	// seed admin
	var userCount int64 = 0
	db.Find(&entity.User{}).Count(&userCount)

	if userCount == 0 {

		pass, err := helpers.PasswordHash("Admin123!@#")
		exception.PanicIfNeeded(err)
		db.Save(&entity.User{
			UserName: "SuperAdmin",
			FullName: "Admin",
			Email:    "admin@mail.com",
			Password: &pass,
		})
	}
}
