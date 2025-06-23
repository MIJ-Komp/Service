package helper

import (
	"fmt"
	"time"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderPaymentServiceHelper struct {
	OrderRepository repository.OrderRepository
}

func NewOrderPaymentServiceHelper(orderRepository repository.OrderRepository) *OrderPaymentServiceHelper {
	return &OrderPaymentServiceHelper{
		OrderRepository: orderRepository,
	}
}

func (service *OrderPaymentServiceHelper) GeneratePaymentRequest(id uuid.UUID, amount float64, desc string, customerId *uuid.UUID, customerEmail string) request.CreatePaymentRequest {
	paymentMethods := []string{"CREDIT_CARD", "BCA", "BNI", "BSI", "BRI", "MANDIRI", "PERMATA", "SAHABAT_SAMPOERNA", "BNC", "ALFAMART", "INDOMARET", "OVO", "DANA", "SHOPEEPAY", "LINKAJA", "JENIUSPAY", "DD_BRI", "DD_BCA_KLIKPAY", "KREDIVO", "AKULAKU", "UANGME", "ATOME", "QRIS"}

	paymentReq := request.CreatePaymentRequest{
		ExternalId:         id.String(),
		Amount:             amount,
		Description:        desc,
		InvoiceDuration:    86400 * 7,
		SuccessRedirectUrl: "http://localhost:3000/order/" + id.String(),
		FailureRedirectUrl: "http://localhost:3000/order/" + id.String(),
		PaymentMethods:     paymentMethods,
		Currency:           "IDR",
		CustomerId:         fmt.Sprintf("%v", customerId),
		CustomerEmail:      customerEmail,
	}

	return paymentReq
}

func (service *OrderPaymentServiceHelper) SetPaidOrder(tx *gorm.DB, paymentId uuid.UUID) {
	order, err := service.OrderRepository.GetByPaymentId(tx, paymentId)
	exception.PanicIfNeeded(err)

	timeNow := time.Now()
	order.IsPaid = true
	order.PaidAt = &timeNow

	_, err = service.OrderRepository.Save(tx, order)
	exception.PanicIfNeeded(err)
}
