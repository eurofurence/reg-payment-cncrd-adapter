package concardis

import (
	"context"
	"time"
)

type Mock interface {
	ConcardisDownstream
}

type mockImpl struct{}

func newMock() Mock {
	return &mockImpl{}
}

func (m *mockImpl) CreatePaymentLink(ctx context.Context, request PaymentLinkCreateRequest) (PaymentLinkCreated, error) {
	return PaymentLinkCreated{
		ID:          42,
		ReferenceID: "deadbeef",
		Link:        "http://localhost:1111/some/paylink",
	}, nil
}

func (m *mockImpl) QueryPaymentLink(ctx context.Context, id uint) (PaymentLinkQueryResponse, error) {
	return PaymentLinkQueryResponse{
		ID:          42,
		Status:      "confirmed",
		ReferenceID: "Order number of my online shop application",
		Link:        "http://localhost:1111/some/paylink",
		Name:        "Online-Shop payment #001",
		Purpose:     "Shop Order #001",
		Amount:      590,
		Currency:    "EUR",
		CreatedAt:   1418392958,
	}, nil
}

func (m *mockImpl) DeletePaymentLink(ctx context.Context, id uint) error {
	return nil
}

func (m *mockImpl) QueryTransactions(ctx context.Context, timeGreaterThan time.Time, timeLessThan time.Time) ([]TransactionData, error) {
	return []TransactionData{}, nil
}
