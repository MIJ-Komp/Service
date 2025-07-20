package database

import (
	"os"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers/logger"
	"api.mijkomp.com/models/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	logger.LogInfo("Connecting to database...")

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})

	if err != nil {
		logger.LogError(err.Error())
		exception.PanicIfNeeded(err)
	}

	logger.LogInfo("Database connected successfully")

	// Auto migrate entity
	db.AutoMigrate(
		&entity.User{},
		&entity.UserVerification{},
		&entity.UserToken{},
		&entity.ProductCategory{},
		&entity.ComponentType{},
		&entity.Product{},
		&entity.ProductSku{},
		&entity.ComponentSpec{},
		&entity.ProductGroupItem{},
		&entity.ProductVariantOption{},
		&entity.ProductVariantOptionValue{},
		&entity.ProductSkuVariant{},
		&entity.VariantOption{},
		&entity.CompatibilityRule{},
		&entity.Menu{},
		&entity.MenuItem{},
		&entity.Order{},
		&entity.Payment{},
		&entity.OrderItem{},
		&entity.CustomerInfo{},
		&entity.ShippingInfo{},
		&entity.Brand{},
	)

	SeedData(db)

	err = CreateSequenceIfNotExists(db, "invoice_code_seq")
	exception.PanicIfNeeded(err)

	return db
}
