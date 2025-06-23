package repository

import (
	"time"

	"api.mijkomp.com/models/entity"
	"api.mijkomp.com/models/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Save(db *gorm.DB, order entity.Order) (entity.Order, error)
	Delete(db *gorm.DB, order entity.Order) error
	Search(db *gorm.DB, query *string, status *enum.EOrderStatus, fromDate, toDate *time.Time, page, pageSize int) ([]entity.Order, int64, int64)
	GetById(db *gorm.DB, orderId *uuid.UUID, code *string) (entity.Order, error)
	GetByPaymentId(db *gorm.DB, paymentId uuid.UUID) (entity.Order, error)
	SaveOrderItems(db *gorm.DB, orderItems []entity.OrderItem) error
	DeleteOrderItems(db *gorm.DB, orderId uuid.UUID, orderItems []entity.OrderItem) error

	// Customer info
	SaveCustomerInfo(db *gorm.DB, customerInfo entity.CustomerInfo) error
	DeleteCustomerInfo(db *gorm.DB, orderId uuid.UUID, customerInfo entity.CustomerInfo) error

	// Shipping info
	SaveShippingInfo(db *gorm.DB, shippingInfo entity.ShippingInfo) error
	DeleteShippingInfo(db *gorm.DB, orderId uuid.UUID, shippingInfo entity.ShippingInfo) error
}
