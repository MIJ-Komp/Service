package service_impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers"
	"api.mijkomp.com/models/entity"
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/repository"
	"api.mijkomp.com/service/helper"

	xendit "github.com/xendit/xendit-go/v4"
	invoice "github.com/xendit/xendit-go/v4/invoice"
	"gorm.io/gorm"
)

type PaymentServiceImpl struct {
	PaymentRepository         repository.PaymentRepository
	OrderPaymentServiceHelper helper.OrderPaymentServiceHelper
	db                        *gorm.DB
}

func NewPaymentService(paymentRepository repository.PaymentRepository, orderRepository repository.OrderRepository, db *gorm.DB) *PaymentServiceImpl {

	orderPaymentServiceHelper := helper.NewOrderPaymentServiceHelper(orderRepository)
	return &PaymentServiceImpl{
		PaymentRepository:         paymentRepository,
		OrderPaymentServiceHelper: *orderPaymentServiceHelper,
		db:                        db,
	}
}

func (service *PaymentServiceImpl) Create(useSqlTransaction bool, request request.CreatePaymentRequest) (*entity.Payment, error) {

	if useSqlTransaction {
		tx := service.db.Begin()
		defer helpers.CommitOrRollback(tx)
	}

	invoiceRequest := *invoice.NewCreateInvoiceRequest(string(request.ExternalId), request.Amount)
	invoiceRequest.Description = &request.Description
	invoiceDuration := fmt.Sprint(request.InvoiceDuration)
	invoiceRequest.InvoiceDuration = &invoiceDuration
	invoiceRequest.SuccessRedirectUrl = &request.SuccessRedirectUrl
	invoiceRequest.FailureRedirectUrl = &request.FailureRedirectUrl
	invoiceRequest.PaymentMethods = request.PaymentMethods
	invoiceRequest.Currency = &request.Currency

	customer := invoice.NewCustomerObject()
	customer.Email = *invoice.NewNullableString(&request.CustomerEmail)
	customer.CustomerId = *invoice.NewNullableString(&request.CustomerId)
	invoiceRequest.Customer = customer

	xenditClient := xendit.NewClient(os.Getenv("XENDIT_SECRET_KEY"))

	resp, r, err := xenditClient.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(invoiceRequest).
		// ForUserId(forUserId). // [OPTIONAL]
		Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InvoiceApi.CreateInvoice``: %v\n", err.Error())

		b, _ := json.Marshal(err.FullError())
		fmt.Fprintf(os.Stderr, "Full Error Struct: %v\n", string(b))

		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		if useSqlTransaction {
			panic(errors.New(string(b)))
		}

		return nil, errors.New(string(b))
	}

	paymentEntity := entity.Payment{
		// CreatedById: userId,
		IsPaid:     false,
		InvoiceUrl: &resp.InvoiceUrl,
		InvoiceId:  resp.Id,
	}

	dbRes, dbErr := service.PaymentRepository.Save(service.db, paymentEntity)

	if dbErr != nil {
		if useSqlTransaction {
			exception.PanicIfNeeded(dbErr)
		}
		return nil, err
	}

	return &dbRes, nil
}

func (service *PaymentServiceImpl) Update(paymentNotification request.PaymentNotification) {
	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	var payment entity.Payment
	err := service.db.First(&payment, "invoice_id = ?", paymentNotification.Id).Error
	exception.PanicIfNeeded(err)

	if paymentNotification.Status == "PAID" {
		payment.IsPaid = true
		payment.TotalPaid = &paymentNotification.PaidAmount
		payment.PaidAt = &paymentNotification.PaidAt
		payment.PayerInfo = &paymentNotification.PayerEmail
		payment.PaymentMethod = &paymentNotification.PaymentMethod
	}

	service.OrderPaymentServiceHelper.SetPaidOrder(tx, payment.Id)

	_, err = service.PaymentRepository.Save(service.db, payment)
	exception.PanicIfNeeded(err)
}
