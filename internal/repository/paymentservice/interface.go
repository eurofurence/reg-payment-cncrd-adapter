package paymentservice

import (
	"context"
	"errors"
	"time"
)

type PaymentService interface {
	AddTransaction(ctx context.Context, transaction Transaction) error
	UpdateTransaction(ctx context.Context, transaction Transaction) error
	GetTransactionByReferenceId(ctx context.Context, reference_id string) (Transaction, error)
}

var (
	NotFoundError   = errors.New("record id not found")
	DownstreamError = errors.New("downstream unavailable - see log for details")
)

type TransactionType string

const (
	Due     TransactionType = "due"
	Payment TransactionType = "payment"
)

type PaymentMethod string

const (
	Credit   PaymentMethod = "credit"
	Cash     PaymentMethod = "cash"
	Paypal   PaymentMethod = "paypal"
	Transfer PaymentMethod = "transfer"
	Internal PaymentMethod = "internal"
	Gift     PaymentMethod = "gift"
)

type TransactionStatus string

const (
	Tentative TransactionStatus = "tentative"
	Pending   TransactionStatus = "pending"
	Valid     TransactionStatus = "valid"
	Deleted   TransactionStatus = "deleted"
)

type StatusHistory struct {
	Status     TransactionStatus `json:"status"`
	Comment    string            `json:"comment"`
	ChangedBy  string            `json:"changed_by"`
	ChangeDate time.Time         `json:"change_date"`
}

type Amount struct {
	Currency  string  `json:"currency"`
	GrossCent int64   `json:"gross_cent"`
	VatRate   float64 `json:"vat_rate"`
}

type Transaction struct {
	DebitorID       uint              `json:"debitor_id"` // TODO this is an 'int64' in the payment service
	ID              string            `json:"transaction_identifier"`
	Type            TransactionType   `json:"transaction_type"`
	Method          PaymentMethod     `json:"method"`
	Amount          Amount            `json:"amount"`
	Comment         string            `json:"comment"`
	Status          TransactionStatus `json:"status"`
	PaymentStartUrl string            `json:"payment_start_url"`
	EffectiveDate   string            `json:"effective_date"`
	DueDate         string            `json:"due_date,omitempty"`
	StatusHistory   []StatusHistory   `json:"status_history"`
}

type TransactionResponse struct {
	Payload []Transaction
}
