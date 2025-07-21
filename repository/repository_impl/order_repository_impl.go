package repository_impl

import (
	"strings"
	"time"

	"api.mijkomp.com/models/entity"
	"api.mijkomp.com/models/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
}

func NewOrderRepository() *OrderRepositoryImpl {
	return &OrderRepositoryImpl{}
}

// Order
func (repository *OrderRepositoryImpl) Save(db *gorm.DB, order entity.Order) (entity.Order, error) {
	var omitFields = []string{}
	err := db.Omit(omitFields...).Save(&order).Error
	return order, err
}

func (repository *OrderRepositoryImpl) Delete(db *gorm.DB, order entity.Order) error {
	err := db.Where("id = ?", order.Id).Update("deleted_at", time.Now().UTC()).Error
	return err
}

func (repository *OrderRepositoryImpl) Search(db *gorm.DB, query *string, status *enum.EOrderStatus, fromDate, toDate *time.Time, page, pageSize int) ([]entity.Order, int64, int64) {
	var orders []entity.Order
	var totalCount int64 = 0
	var totalPage int64 = 0
	var offset int = 0

	queries := db.Model(&orders).
		Preload("OrderItems", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
		Preload("OrderItems.Product").
		Preload("OrderItems.Product.ProductSkus", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
		Preload("CustomerInfo").
		Preload("ShippingInfo").
		Preload("Payment")

	if query != nil {
		queries.Where("name like ?", "%"+*query+"%")
	}

	if status != nil {
		queries.Where("status = ?", status)
	}

	if fromDate != nil {
		queries.Where("order_date >== ?", fromDate)
	}

	if toDate != nil {
		queries.Where("order_date <= ?", toDate)
	}

	queries.Count(&totalCount)

	offset = (page - 1) * pageSize
	totalPage = ((totalCount + int64(pageSize) - 1) / int64(pageSize))

	queries.Limit(pageSize).Offset(offset).Order("orders.order_date desc").Find(&orders)

	return orders, totalCount, totalPage
}

func (repository *OrderRepositoryImpl) GetById(db *gorm.DB, orderId *uuid.UUID, code *string) (entity.Order, error) {
	var order entity.Order
	queries := db.Model(&entity.Order{}).
		Preload("OrderItems", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
		Preload("OrderItems.Product").
		Preload("OrderItems.Product.ProductSkus").
		Preload("CustomerInfo").
		Preload("ShippingInfo").
		Preload("Payment")

	if orderId != nil {
		queries.Where("id = ?", orderId)
	}

	if code != nil {
		queries.Where("LOWER(code) = ?", strings.ToLower(*code))
	}

	err := queries.First(&order).Error
	return order, err
}

func (repository *OrderRepositoryImpl) GetByPaymentId(db *gorm.DB, paymentId uuid.UUID) (entity.Order, error) {
	var order entity.Order
	err := db.
		Preload("OrderItems", func(db *gorm.DB) *gorm.DB { return db.Order("sequence") }).
		Preload("OrderItems.Product").
		Preload("OrderItems.Product.ProductSkus").
		Preload("CustomerInfo").
		Preload("ShippingInfo").
		First(&order, "payment_id = ?", paymentId).Error
	return order, err
}

// Order Items
func (repository *OrderRepositoryImpl) SaveOrderItems(db *gorm.DB, orderItems []entity.OrderItem) error {
	err := db.Save(&orderItems).Error
	return err
}

func (repository *OrderRepositoryImpl) DeleteOrderItems(db *gorm.DB, orderId uuid.UUID, orderItems []entity.OrderItem) error {
	err := db.Where("order_id = ?", orderId).Delete(orderItems).Error
	return err
}

// Customer info
func (repository *OrderRepositoryImpl) SaveCustomerInfo(db *gorm.DB, customerInfo entity.CustomerInfo) error {
	err := db.Save(&customerInfo).Error
	return err
}

func (repository *OrderRepositoryImpl) DeleteCustomerInfo(db *gorm.DB, orderId uuid.UUID, customerInfo entity.CustomerInfo) error {
	err := db.Where("order_id = ?", orderId).Delete(customerInfo).Error
	return err
}

// Shipping info
func (repository *OrderRepositoryImpl) SaveShippingInfo(db *gorm.DB, shippingInfo entity.ShippingInfo) error {
	err := db.Save(&shippingInfo).Error
	return err
}

func (repository *OrderRepositoryImpl) DeleteShippingInfo(db *gorm.DB, orderId uuid.UUID, shippingInfo entity.ShippingInfo) error {
	err := db.Where("order_id = ?", orderId).Delete(shippingInfo).Error
	return err
}
