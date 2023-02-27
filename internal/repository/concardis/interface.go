package concardis

import (
	"context"
	"errors"
	"time"
)

type ConcardisDownstream interface {
	CreatePaymentLink(ctx context.Context, request PaymentLinkCreateRequest) (PaymentLinkCreated, error)
	QueryPaymentLink(ctx context.Context, id uint) (PaymentLinkQueryResponse, error)
	DeletePaymentLink(ctx context.Context, id uint) error

	QueryTransactions(ctx context.Context, timeGreaterThan time.Time, timeLessThan time.Time) ([]TransactionData, error)
}

var (
	NoSuchID404Error = errors.New("payment link id not found")
	DownstreamError  = errors.New("downstream unavailable - see log for details")
	NotSuccessful    = errors.New("response body status field did not indicate success")
)

// -- CreatePaymentLink --

type PaymentLinkCreateRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	PSP         uint64  `json:"psp"`
	ReferenceId string  `json:"referenceId"`
	OrderId     string  `json:"concardisOrderId"`
	Purpose     string  `json:"purpose"`
	Amount      int64   `json:"amount"`  // in cents
	VatRate     float64 `json:"vatRate"` // in %
	Currency    string  `json:"currency"`
	SKU         string  `json:"sku"`
}

type PaymentLinkCreated struct {
	ID          uint   `json:"id"`
	ReferenceID string `json:"referenceId"`
	Link        string `json:"link"`
}

// -- QueryPaymentLink --

type PaymentLinkQueryResponse struct {
	ID          uint                 `json:"id"` // not the payment link id!
	Status      string               `json:"status"`
	ReferenceID string               `json:"referenceId"`
	Link        string               `json:"link"`
	Invoices    []PaymentLinkInvoice `json:"invoices"`
	Name        string               `json:"name"`
	Purpose     map[string]string    `json:"purpose"`
	Amount      int64                `json:"amount"`
	Currency    string               `json:"currency"`
	CreatedAt   int64                `json:"createdAt"`
}

type PaymentLinkInvoice struct {
	ReferenceID      string            `json:"referenceId"`
	PaymentRequestId uint              `json:"paymentRequestId"` // the payment link id
	Currency         string            `json:"currency"`
	Amount           int64             `json:"amount"`
	Transactions     []TransactionData `json:"transactions"`
}

// -- QueryTransactions --

type TransactionData struct {
	ID          int64   `json:"id"`
	UUID        string  `json:"uuid"`
	Amount      int64   `json:"amount"`
	Status      string  `json:"status"`
	Time        string  `json:"time"`
	Lang        string  `json:"lang"`
	PageUUID    string  `json:"pageUuid"`
	Payment     Payment `json:"payment"`
	Psp         string  `json:"psp"`
	PspID       int64   `json:"pspId"`
	Mode        string  `json:"mode"` // "LIVE"
	ReferenceID string  `json:"referenceId"`
	Invoice     Invoice `json:"invoice"`
}

type Payment struct {
	Brand string `json:"brand"`
}

type Invoice struct {
	ReferenceID      string `json:"referenceId"`
	PaymentRequestId uint   `json:"paymentRequestId"` // the payment link id
	Currency         string `json:"currency"`
	OriginalAmount   int64  `json:"originalAmount"`
	RefundedAmount   int64  `json:"refundedAmount"`
}
