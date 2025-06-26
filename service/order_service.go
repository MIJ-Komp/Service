package service

import (
	"time"

	"api.mijkomp.com/models/enum"
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
	"github.com/google/uuid"
)

type OrderService interface {
	Create(payload request.Order) response.Order
	Update(orderId uuid.UUID, payload request.Order) response.Order
	UpdateStatus(currentUserId uint, orderId uuid.UUID, payload request.UpdateOrderStatusByAdmin) response.Order
	UpdateShippingInfo(currentUserId uint, orderId uuid.UUID, payload request.UpdateOrderShippingByAdmin) response.Order
	Delete(orderId uuid.UUID) string
	Search(query *string, status *enum.EOrderStatus, fromDate, toDate *time.Time, page, pageSize int) response.PageResult
	GetById(orderId *uuid.UUID, code *string) response.Order
}
