package database

import (
	"os"

	"api.mijkomp.com/exception"
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
		&entity.ProductSpec{},
		&entity.ProductGroupItem{},
		&entity.ProductVariantOption{},
		&entity.ProductVariantOptionValue{},
		&entity.ProductSkuVariant{},
		&entity.VariantOption{},
		&entity.ComponentType{},
		&entity.CompatibilityRule{},
		&entity.Menu{},
		&entity.MenuItem{},
		&entity.Order{},
		&entity.Payment{},
		&entity.OrderItem{},
		&entity.CustomerInfo{},
		&entity.ShippingInfo{},
	)

	// SeedData(db)

	// err = CreateSequenceIfNotExists(db, "invoice_code_seq")
	// exception.PanicIfNeeded(err)

	return db
}
