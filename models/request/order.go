package request

import (
	"time"

	"api.mijkomp.com/models/enum"
	"github.com/google/uuid"
)

type Order struct {
	Notes string `json:"notes"`

	CustomerInfo CustomerInfo `json:"customerInfo"`
	OrderItems   []OrderItem  `json:"orderItems"`
	ShippingInfo ShippingInfo `json:"shippingInfo"`
}

type OrderItem struct {
	ProductId    uuid.UUID `json:"productId"`
	ProductSkuId uuid.UUID `json:"productSkuId"`
	Quantity     int       `json:"quantity"`
}

type CustomerInfo struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

type ShippingInfo struct {
	RecipientName string  `json:"recipientName"`
	Address       string  `json:"address"`
	Province      string  `json:"province"`
	City          string  `json:"city"`
	PostalCode    string  `json:"postalCode"`
	Notes         *string `json:"notes"`
}

type UpdateOrderStatusByAdmin struct {
	NewStatus enum.EOrderStatus `json:"newStatus"`
}

type UpdateOrderShippingByAdmin struct {
	ShippingMethod    *string    `json:"shippingMethod"`
	TrackingNumber    *string    `json:"trackingNumber"`
	EstimatedDelivery *time.Time `json:"estimatedDelivery"`
	ShippedAt         *time.Time `json:"shippedAt"`
	DeliveredAt       *time.Time `json:"deliveredAt"`
}
