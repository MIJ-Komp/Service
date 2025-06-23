package service

import (
	"api.mijkomp.com/models/entity"
	"api.mijkomp.com/models/request"
)

type PaymentService interface {
	Create(useSqlTransaction bool, request request.CreatePaymentRequest) (*entity.Payment, error)
	Update(paymentNotification request.PaymentNotification)
}
