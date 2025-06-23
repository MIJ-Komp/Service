package response

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	Id         uuid.UUID    `json:"id"`
	Code       string       `json:"code"`
	CustomerId *uint        `json:"customerId"`
	OrderDate  time.Time    `json:"orderDate"`
	Status     EnumResponse `json:"status"`
	Notes      string       `json:"notes"`

	IsPaid     bool       `json:"isPaid"`
	TotalPaid  *float64   `json:"totalPaid"`
	PaidAt     *time.Time `json:"paidAt"`
	PaymentUrl string     `json:"paymentUrl"`

	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`

	OrderItems   []OrderItem  `json:"orderItems"`
	CustomerInfo CustomerInfo `json:"customerInfo"`
	ShippingInfo ShippingInfo `json:"shippingInfo"`
}

type OrderItem struct {
	Id           uuid.UUID `json:"id"`
	OrderId      uuid.UUID `json:"orderId"`
	ProductId    uuid.UUID `json:"productId"`
	ProductSkuId uuid.UUID `json:"productSkuId"`
	Quantity     int       `json:"quantity"`
	Price        float64   `json:"price"`

	Product OrderItemProduct `json:"product"`
}

type OrderItemProduct struct {
	Id          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	SKU         string      `json:"sku"`
	IsActive    bool        `json:"isActive"`
	ImageIds    []uuid.UUID `json:"imageIds"`
	Description string      `json:"description"`

	ProductSku OrderItemProductSku `json:"productSku"`
}

type OrderItemProductSku struct {
	Id    uuid.UUID `json:"id"`
	SKU   string    `json:"sku"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
	Stock *int      `json:"stock"`
}

type CustomerInfo struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
}

type ShippingInfo struct {
	Id                uuid.UUID  `json:"id"`
	OrderId           uuid.UUID  `json:"orderId"`
	RecipientName     string     `json:"recipientName"`
	Address           string     `json:"address"`
	Province          string     `json:"province"`
	City              string     `json:"city"`
	PostalCode        string     `json:"postalCode"`
	ShippingMethod    *string    `json:"shippingMethod"`
	TrackingNumber    *string    `json:"trackingNumber"`
	EstimatedDelivery *time.Time `json:"estimatedDelivery"`
	ShippedAt         *time.Time `json:"shippedAt"`
	DeliveredAt       *time.Time `json:"deliveredAt"`
	Notes             *string    `json:"notes"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}
