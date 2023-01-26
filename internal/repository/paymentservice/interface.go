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

type Deletion struct {
	PreviousStatus TransactionStatus
	Comment        string
	DeletedBy      string
	Date           time.Time
}

type Amount struct {
	Currency  string
	GrossCent int64
	VatRate   float64
}

type Transaction struct {
	ID              string            `json:"transaction_identifier"`
	DebitorID       uint              `json:"debitor_id"` // XXX TODO this is an 'int64' in the payment service
	Type            TransactionType   `json:"transaction_type"`
	Method          PaymentMethod     `json:"method"`
	Amount          Amount            `json:"amount"`
	Comment         string            `json:"comment"`
	Status          TransactionStatus `json:"status"`
	PaymentStartUrl string            `json:"payment_start_url"`
	EffectiveDate   string            `json:"effective_date"`
	DueDate         string            `json:"due_date,omitempty"`
	Deletion        *Deletion
}

type TransactionResponse struct {
	Payload []Transaction
}
