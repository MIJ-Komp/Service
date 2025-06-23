package request

import "time"

type CreatePaymentRequest struct {
	ExternalId         string
	Amount             float64
	Description        string
	InvoiceDuration    uint
	SuccessRedirectUrl string
	FailureRedirectUrl string
	PaymentMethods     []string
	Currency           string
	CustomerId         string
	CustomerEmail      string
}

type PaymentNotification struct {
	Id                 string    `json:"id"`                  // "593f4ed1c3d3bb7f39733d83",
	ExternalId         string    `json:"external_id"`         // "testing-invoice",
	UserId             string    `json:"user_id"`             // "5848fdf860053555135587e7",
	IsHigh             bool      `json:"is_high"`             // false,
	PaymentMethod      string    `json:"payment_method"`      // "BANK_TRANSFER",
	Status             string    `json:"status"`              // "PAID",
	MerchantName       string    `json:"merchant_name"`       // "Xendit",
	Amount             float64   `json:"amount"`              // 2000000,
	PaidAmount         float64   `json:"paid_amount"`         // 2000000,
	BankCode           string    `json:"bank_code"`           // "MANDIRI",
	PaidAt             time.Time `json:"paid_at"`             // "2020-01-14T02:32:50.912Z",
	PayerEmail         string    `json:"payer_email"`         // "test@xendit.co",
	Description        string    `json:"description"`         // "Invoice webhook test",
	Created            time.Time `json:"created"`             // "2020-01-13T02:32:49.827Z",
	Updated            time.Time `json:"updated"`             // "2020-01-13T02:32:50.912Z",
	Currency           string    `json:"currency"`            // "IDR",
	PaymentChannel     string    `json:"payment_channel"`     // "MANDIRI",
	PaymentDestination string    `json:"payment_destination"` // "8458478548758748"
}
