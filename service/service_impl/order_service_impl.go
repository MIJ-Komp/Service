package service_impl

import (
	"fmt"
	"slices"
	"time"

	"api.mijkomp.com/database"
	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers"
	"api.mijkomp.com/models/entity"
	"api.mijkomp.com/models/enum"
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/repository"
	"api.mijkomp.com/service"
	"api.mijkomp.com/service/helper"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderServiceImpl struct {
	OrderRepository           repository.OrderRepository
	ProductService            service.ProductService
	ProductRepository         repository.ProductRepository
	OrderPaymentServiceHelper helper.OrderPaymentServiceHelper
	PaymentService            service.PaymentService
	db                        *gorm.DB
}

func NewOrderService(orderRepostitory repository.OrderRepository, productService service.ProductService, productRepository repository.ProductRepository, paymentService service.PaymentService, db *gorm.DB) *OrderServiceImpl {
	orderPaymentServiceHelper := helper.NewOrderPaymentServiceHelper(orderRepostitory)
	return &OrderServiceImpl{
		OrderRepository:           orderRepostitory,
		ProductService:            productService,
		ProductRepository:         productRepository,
		OrderPaymentServiceHelper: *orderPaymentServiceHelper,
		PaymentService:            paymentService,
		db:                        db,
	}
}

func (service *OrderServiceImpl) Create(payload request.Order) response.Order {

	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	newId := uuid.New()
	newCode, err := database.GetNextInvoiceNumber(tx)
	exception.PanicIfNeeded(err)

	formattedCode := fmt.Sprintf("INV%05d", newCode)

	order := entity.Order{
		Id:                   newId,
		Code:                 formattedCode,
		OrderDate:            time.Now(),
		Status:               enum.OrderStatusPending,
		Notes:                payload.Notes,
		CreatedByCustomerAt:  time.Now(),
		ModifiedByCustomerAt: time.Now(),
	}

	// Order Items

	productSkuIds := []uuid.UUID{}
	for _, orderItem := range payload.OrderItems {
		productSkuIds = append(productSkuIds, orderItem.ProductSkuId)
	}

	productSkus := service.ProductRepository.GetProductSkuByIds(tx, productSkuIds)
	orderItems := []entity.OrderItem{}

	var totalPrice float64 = 0
	for i, orderItem := range payload.OrderItems {

		productSkuIdx := slices.IndexFunc(productSkus, func(model response.BrowseProductSku) bool {
			return model.Id == orderItem.ProductSkuId && model.ProductId == orderItem.ProductId
		})

		if productSkuIdx == -1 {
			panic(exception.NewValidationError(fmt.Sprintf("Product at index %d not found.", i)))
		}

		productSku := productSkus[productSkuIdx]

		if productSku.Stock != nil {
			if *productSku.Stock < orderItem.Quantity {
				partMsg := "telah habis"
				if *productSku.Stock > 0 {
					partMsg = fmt.Sprintf("tersisa %d", *productSku.Stock)
				}

				panic(exception.NewValidationError(fmt.Sprintf("Stok produk dengan sku %s %s", productSku.SKU, partMsg)))
			}

			// kurangi qty product
			*productSkus[productSkuIdx].Stock -= orderItem.Quantity
		}

		orderItems = append(orderItems, entity.OrderItem{
			Id:           uuid.New(),
			OrderId:      newId,
			ProductId:    orderItem.ProductId,
			ProductSkuId: orderItem.ProductSkuId,
			Quantity:     orderItem.Quantity,
			Price:        productSku.Price,
			Sequence:     i + 1,
			CreatedAt:    time.Now(),
		})

		totalPrice += productSku.Price * float64(orderItem.Quantity)
	}

	// create payment xendit
	desc := "Order Payment MIJKomp"
	paymentReq := service.OrderPaymentServiceHelper.GeneratePaymentRequest(newId, totalPrice, desc, nil, payload.CustomerInfo.Email)

	paymentInfo, err := service.PaymentService.Create(false, paymentReq)
	exception.PanicIfNeeded(err)

	order.PaymentId = &paymentInfo.Id

	// Save Order
	newOrder, err := service.OrderRepository.Save(tx, order)
	exception.PanicIfNeeded(err)

	err = service.OrderRepository.SaveOrderItems(tx, orderItems)
	exception.PanicIfNeeded(err)

	// Customer Info
	customerInfo := entity.CustomerInfo{
		OrderId:     newOrder.Id,
		Name:        payload.CustomerInfo.Name,
		Email:       payload.CustomerInfo.Email,
		PhoneNumber: payload.CustomerInfo.PhoneNumber,
	}

	err = service.OrderRepository.SaveCustomerInfo(tx, customerInfo)
	exception.PanicIfNeeded(err)

	// Shipping Info
	shippingInfo := entity.ShippingInfo{
		Id:            uuid.New(),
		OrderId:       newOrder.Id,
		RecipientName: payload.ShippingInfo.RecipientName,
		Address:       payload.ShippingInfo.Address,
		District:      payload.ShippingInfo.District,
		City:          payload.ShippingInfo.City,
		Province:      payload.ShippingInfo.Province,
		PostalCode:    payload.ShippingInfo.PostalCode,
		Notes:         payload.ShippingInfo.Notes,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	err = service.OrderRepository.SaveShippingInfo(tx, shippingInfo)
	exception.PanicIfNeeded(err)

	// Update stock product sku
	productSkusToBeUpdated := []entity.ProductSku{}
	for _, pSku := range productSkus {
		if pSku.ProductType != string(enum.ProductTypeGroup) {
			productSkusToBeUpdated = append(productSkusToBeUpdated, entity.ProductSku{
				Id:    pSku.Id,
				Stock: pSku.Stock,
			})
		}
	}

	err = service.ProductRepository.UpdateStockProductSkus(tx, productSkusToBeUpdated)
	exception.PanicIfNeeded(err)

	createdOrder, err := service.OrderRepository.GetById(tx, &newOrder.Id, nil)
	exception.PanicIfNeeded(err)
	return service.mapOrder(createdOrder)
}

func (service *OrderServiceImpl) Update(orderId uuid.UUID, payload request.Order) response.Order {
	// tx
	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	order, err := service.OrderRepository.GetById(service.db, &orderId, nil)
	exception.PanicIfNeeded(err)

	order.Notes = payload.Notes
	order.ModifiedByCustomerAt = time.Now().UTC()

	_, err = service.OrderRepository.Save(tx, order)
	exception.PanicIfNeeded(err)

	// Order Items (Update, Delete, Add New)
	orderItems := order.OrderItems
	orderItemToBeDeleted := []entity.OrderItem{}

	// Update, delete
	for i, orderItem := range orderItems {
		payloadIdx := slices.IndexFunc(payload.OrderItems, func(model request.OrderItem) bool {
			return model.ProductSkuId == orderItem.ProductSkuId
		})

		if payloadIdx != -1 {
			orderItems[i].Quantity = payload.OrderItems[payloadIdx].Quantity
			orderItems[i].Price = 5000
		} else {
			orderItemToBeDeleted = append(orderItemToBeDeleted, orderItem)
		}
	}

	// Add new
	for _, orderItem := range payload.OrderItems {
		savedIdx := slices.IndexFunc(orderItems, func(model entity.OrderItem) bool {
			return model.ProductSkuId == orderItem.ProductSkuId
		})

		if savedIdx == -1 {
			orderItems = append(orderItems, entity.OrderItem{
				OrderId:      order.Id,
				ProductSkuId: orderItem.ProductSkuId,
				Price:        5000,
				Quantity:     orderItem.Quantity,
				Sequence:     len(orderItems) + 1,
				CreatedAt:    time.Now().UTC(),
			})
		}
	}

	err = service.OrderRepository.DeleteOrderItems(tx, order.Id, orderItemToBeDeleted)
	exception.PanicIfNeeded(err)

	err = service.OrderRepository.SaveOrderItems(tx, orderItems)
	exception.PanicIfNeeded(err)

	// return new result
	res, err := service.OrderRepository.GetById(tx, &orderId, nil)
	exception.PanicIfNeeded(err)

	return service.mapOrder(res)
}

func (service *OrderServiceImpl) UpdateStatus(currentUserId uint, orderId uuid.UUID, payload request.UpdateOrderStatusByAdmin) response.Order {
	// tx
	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	order, err := service.OrderRepository.GetById(service.db, &orderId, nil)
	exception.PanicIfNeeded(err)

	order.Status = payload.NewStatus
	order.ModifiedByAdminId = &currentUserId
	timeNow := time.Now()
	order.ModifiedByAdminAt = &timeNow

	_, err = service.OrderRepository.Save(tx, order)
	exception.PanicIfNeeded(err)

	// return new result
	res, err := service.OrderRepository.GetById(tx, &orderId, nil)
	exception.PanicIfNeeded(err)

	return service.mapOrder(res)
}

func (service *OrderServiceImpl) UpdateShippingInfo(currentUserId uint, orderId uuid.UUID, payload request.UpdateOrderShippingByAdmin) response.Order {
	// tx
	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	order, err := service.OrderRepository.GetById(service.db, &orderId, nil)
	exception.PanicIfNeeded(err)

	shippingInfo := order.ShippingInfo

	shippingInfo.ShippingMethod = payload.ShippingMethod
	shippingInfo.TrackingNumber = payload.TrackingNumber
	shippingInfo.EstimatedDelivery = payload.EstimatedDelivery
	shippingInfo.ShippedAt = payload.ShippedAt
	shippingInfo.DeliveredAt = payload.DeliveredAt

	order.ModifiedByAdminId = &currentUserId
	timeNow := time.Now()
	order.ModifiedByAdminAt = &timeNow

	err = service.OrderRepository.SaveShippingInfo(tx, shippingInfo)
	exception.PanicIfNeeded(err)

	// return new result
	res, err := service.OrderRepository.GetById(tx, &orderId, nil)
	exception.PanicIfNeeded(err)

	return service.mapOrder(res)
}

func (service *OrderServiceImpl) Delete(orderId uuid.UUID) string {
	// tx
	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	order, err := service.OrderRepository.GetById(service.db, &orderId, nil)
	exception.PanicIfNeeded(err)

	err = service.OrderRepository.Delete(service.db, order)
	exception.PanicIfNeeded(err)

	return fmt.Sprintf("Order %s berhasil di hapus", order.Code)
}

func (service *OrderServiceImpl) Search(query *string, status *enum.EOrderStatus, fromDate, toDate *time.Time, page, pageSize int) response.PageResult {

	res, totalCount, totalPage := service.OrderRepository.Search(service.db, query, status, fromDate, toDate, page, pageSize)

	return response.PageResult{
		Items:      service.mapOrders(res),
		TotalCount: totalCount,
		PageSize:   totalPage,
	}
}

func (service *OrderServiceImpl) GetById(orderId *uuid.UUID, code *string) response.Order {
	res, err := service.OrderRepository.GetById(service.db, orderId, code)
	exception.PanicIfNeeded(err)

	return service.mapOrder(res)
}

// Map helpers
func (service *OrderServiceImpl) mapOrders(orders []entity.Order) []response.Order {
	orderRes := []response.Order{}

	for _, el := range orders {
		orderRes = append(orderRes, service.mapOrder(el))
	}

	return orderRes
}

func (service *OrderServiceImpl) mapOrder(order entity.Order) response.Order {
	orderRes := response.Order{
		Id:         order.Id,
		Code:       order.Code,
		CustomerId: order.CustomerId,
		OrderDate:  order.OrderDate,
		Status: response.EnumResponse{
			Code: string(order.Status),
			Name: order.Status.DisplayString(),
		},
		PaidAt:     order.Payment.PaidAt,
		IsPaid:     order.IsPaid,
		TotalPaid:  order.Payment.TotalPaid,
		PaymentUrl: order.Payment.InvoiceUrl,
		Notes:      order.Notes,

		CreatedByCustomerAt:  order.CreatedByCustomerAt,
		ModifiedByCustomerAt: order.ModifiedByCustomerAt,

		OrderItems:   service.mapOrderItems(order.OrderItems),
		CustomerInfo: service.mapCustomerInfo(order.CustomerInfo),
		ShippingInfo: service.mapShippingInfo(order.ShippingInfo),
	}

	return orderRes
}

// Map Order Items

func (service *OrderServiceImpl) mapOrderItems(orderItems []entity.OrderItem) []response.OrderItem {
	res := []response.OrderItem{}

	for _, item := range orderItems {
		res = append(res, response.OrderItem{
			Id:           item.Id,
			OrderId:      item.OrderId,
			ProductId:    item.ProductId,
			ProductSkuId: item.ProductSkuId,
			Quantity:     item.Quantity,
			Price:        item.Price,

			Product: service.mapOrderItemProduct(item.Product, item.ProductSkuId),
		})
	}

	return res
}

func (service *OrderServiceImpl) mapOrderItemProduct(product entity.Product, productSkuId uuid.UUID) response.OrderItemProduct {
	res := response.OrderItemProduct{
		Id:          product.Id,
		Name:        product.Name,
		SKU:         product.SKU,
		IsActive:    product.IsActive,
		ImageIds:    helpers.SplitImageIds(product.ImageIds),
		Description: product.Description,
	}

	productSkuIdx := slices.IndexFunc(product.ProductSkus, func(model entity.ProductSku) bool {
		return model.Id == productSkuId
	})

	if productSkuIdx != -1 {
		res.ProductSku = response.OrderItemProductSku{
			Id:    product.ProductSkus[productSkuIdx].Id,
			SKU:   product.ProductSkus[productSkuIdx].SKU,
			Name:  product.ProductSkus[productSkuIdx].Name,
			Price: product.ProductSkus[productSkuIdx].Price,
			Stock: product.ProductSkus[productSkuIdx].Stock,
		}
	}

	return res

}

func (service *OrderServiceImpl) mapCustomerInfo(customerInfo entity.CustomerInfo) response.CustomerInfo {
	return response.CustomerInfo{
		Id:          customerInfo.Id,
		Name:        customerInfo.Name,
		Email:       customerInfo.Email,
		PhoneNumber: customerInfo.PhoneNumber,
	}
}

func (service *OrderServiceImpl) mapShippingInfo(shippingInfo entity.ShippingInfo) response.ShippingInfo {
	return response.ShippingInfo{
		Id:                shippingInfo.Id,
		OrderId:           shippingInfo.OrderId,
		RecipientName:     shippingInfo.RecipientName,
		Address:           shippingInfo.Address,
		District:          shippingInfo.District,
		City:              shippingInfo.City,
		Province:          shippingInfo.Province,
		PostalCode:        shippingInfo.PostalCode,
		ShippingMethod:    shippingInfo.ShippingMethod,
		TrackingNumber:    shippingInfo.TrackingNumber,
		EstimatedDelivery: shippingInfo.EstimatedDelivery,
		ShippedAt:         shippingInfo.ShippedAt,
		DeliveredAt:       shippingInfo.DeliveredAt,
		Notes:             shippingInfo.Notes,
		CreatedAt:         shippingInfo.CreatedAt,
		UpdatedAt:         shippingInfo.UpdatedAt,
	}
}
