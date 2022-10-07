package concardis

import (
	"context"
	"fmt"
	"time"
)

type Mock interface {
	ConcardisDownstream

	Reset()
	Recording() []string
	SimulateError(err error)
}

type mockImpl struct {
	recording     []string
	simulateError error
}

func newMock() Mock {
	return &mockImpl{
		recording: make([]string, 0),
	}
}

func (m *mockImpl) CreatePaymentLink(ctx context.Context, request PaymentLinkCreateRequest) (PaymentLinkCreated, error) {
	if m.simulateError != nil {
		return PaymentLinkCreated{}, m.simulateError
	}
	m.recording = append(m.recording, fmt.Sprintf("CreatePaymentLink %v", request))
	return PaymentLinkCreated{
		ID:          42,
		ReferenceID: "deadbeef",
		Link:        "http://localhost:1111/some/paylink",
	}, nil
}

func (m *mockImpl) QueryPaymentLink(ctx context.Context, id uint) (PaymentLinkQueryResponse, error) {
	if m.simulateError != nil {
		return PaymentLinkQueryResponse{}, m.simulateError
	}
	m.recording = append(m.recording, fmt.Sprintf("QueryPaymentLink %d", id))
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
	if m.simulateError != nil {
		return m.simulateError
	}
	m.recording = append(m.recording, fmt.Sprintf("DeletePaymentLink %d", id))
	return nil
}

func (m *mockImpl) QueryTransactions(ctx context.Context, timeGreaterThan time.Time, timeLessThan time.Time) ([]TransactionData, error) {
	if m.simulateError != nil {
		return []TransactionData{}, m.simulateError
	}
	m.recording = append(m.recording, fmt.Sprintf("QueryTransactions %v <= t <= %v", timeGreaterThan, timeLessThan))
	return []TransactionData{}, nil
}

func (m *mockImpl) Reset() {
	m.recording = make([]string, 0)
	m.simulateError = nil
}

func (m *mockImpl) Recording() []string {
	return m.recording
}

func (m *mockImpl) SimulateError(err error) {
	m.simulateError = err
}
