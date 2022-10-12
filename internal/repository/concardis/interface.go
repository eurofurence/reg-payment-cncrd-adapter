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
)

// -- CreatePaymentLink --

type PaymentLinkCreateRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	PSP         uint64  `json:"psp"`
	ReferenceId string  `json:"referenceId"`
	Purpose     string  `json:"purpose"`
	Amount      int64   `json:"amount"` // TODO is this cents?
	VatRate     float64 `json:"vatRate"`
	Currency    string  `json:"currency"`
	SKU         string  `json:"sku"`
}

type PaymentLinkCreated struct {
	ID          int64  `json:"id"`
	ReferenceID string `json:"referenceId"`
	Link        string `json:"link"`
}

// -- QueryPaymentLink --

type PaymentLinkQueryResponse struct {
	ID          int64  `json:"id"`
	Status      string `json:"status"`
	ReferenceID string `json:"referenceId"`
	Link        string `json:"link"`
	Name        string `json:"name"`
	Purpose     string `json:"purpose"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	CreatedAt   int64  `json:"createdAt"`
}

// -- QueryTransactions --

type TransactionData struct {
	ID          int64   `json:"id"`
	UUID        string  `json:"uuid"`
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
	CurrencyAlpha3 string `json:"currencyAlpha3"`
	ShippingAmount int64  `json:"shippingAmount"`
	TotalAmount    int64  `json:"totalAmount"`
}
