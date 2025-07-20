package service_impl

import (
	"time"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/models/entity"
	"api.mijkomp.com/models/enum"
	"api.mijkomp.com/models/response"
	"gorm.io/gorm"
)

type DashboardServiceImpl struct {
	db *gorm.DB
}

func NewDashboardService(db *gorm.DB) *DashboardServiceImpl {
	return &DashboardServiceImpl{
		db: db,
	}
}

func (service *DashboardServiceImpl) GetSummary(currentUserId uint, fromDate, toDate time.Time) response.Dashboard {

	var totalOrder int64
	var totalPendingOrder int64
	var totalSales float64
	var totalActiveProduct int64

	stockAlert := []response.StockAlert{}
	bestSellingProduct := []response.BestSellingProduct{}

	// total sales
	service.db.Model(&entity.Order{}).
		Select("COALESCE(SUM(total_paid), 0)").
		Where("is_paid = true AND total_paid IS NOT NULL AND deleted_at IS NULL AND status = ? AND created_at >= ? AND created_at <= ?", enum.OrderStatusCompleted, fromDate, toDate).
		Scan(&totalSales)

	// total order
	service.db.Model(&entity.Order{}).
		Where("deleted_at IS NULL and status != ? AND created_at >= ? AND created_at <= ?", enum.OrderStatusCancelled, fromDate, toDate).
		Count(&totalOrder)

	// pending order
	service.db.Model(&entity.Order{}).
		Where("deleted_at IS NULL and status not in ? AND created_at >= ? AND created_at <= ?", []enum.EOrderStatus{enum.OrderStatusCompleted, enum.OrderStatusCancelled}, fromDate, toDate).
		Count(&totalPendingOrder)

	// active product
	service.db.
		Table("product_skus").
		Joins("JOIN products ON products.id = product_skus.product_id").
		Where("product_skus.is_active = true").
		Where("product_skus.stock > 0").
		Where("products.deleted_at IS NULL").
		Where("products.is_active = true AND products.deleted_at IS NULL").
		Count(&totalActiveProduct)

	// stock alert
	err := service.db.
		Table("products").
		Joins("JOIN product_skus ON products.id = product_skus.product_id").
		Where("product_skus.is_active = true").
		Where("product_skus.stock <= product_skus.stock_alert").
		Where("product_type = ?", enum.ProductTypeSingle).
		Where("products.is_active = true AND products.deleted_at IS NULL").
		Limit(10).
		Select("products.id, products.name || product_skus.name AS name, product_skus.stock").
		Scan(&stockAlert).Error
	exception.PanicIfNeeded(err)

	// Best selling product
	service.db.Table("order_items AS oi").
		Select("p.id AS id, p.name AS name, SUM(oi.quantity) AS sold").
		Joins("JOIN product_skus ps ON ps.id = oi.product_sku_id").
		Joins("JOIN products p ON p.id = ps.product_id").
		Joins("JOIN orders o ON o.id = oi.order_id").
		Where("o.is_paid = true AND o.deleted_at IS NULL AND p.deleted_at IS NULL").
		Group("p.id, p.name").
		Order("sold DESC").
		Limit(10).
		Scan(&bestSellingProduct)

	return response.Dashboard{
		TotalSales:         totalSales,
		TotalOrder:         totalOrder,
		TotalPendingOrder:  totalPendingOrder,
		TotalActiveProduct: totalActiveProduct,
		StockAlert:         stockAlert,
		BestSellingProduct: bestSellingProduct,
	}
}
