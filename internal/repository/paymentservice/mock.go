package paymentservice

import (
	"context"
)

type Mock interface {
	PaymentService

	InjectTransaction(ctx context.Context, transaction Transaction) error
	Reset()
	Recording() []Transaction
	SimulateAddError(err error)
}

type MockImpl struct {
	data             map[uint][]Transaction
	recording        []Transaction
	simulateGetError error
	simulateAddError error
}

var (
	_ PaymentService = (*MockImpl)(nil)
	_ Mock           = (*MockImpl)(nil)
)

func newMock() Mock {
	return &MockImpl{
		data:      make(map[uint][]Transaction),
		recording: make([]Transaction, 0),
	}
}

func (m *MockImpl) AddTransaction(ctx context.Context, transaction Transaction) error {
	if m.simulateAddError != nil {
		return m.simulateAddError
	}

	_ = m.InjectTransaction(ctx, transaction)
	m.recording = append(m.recording, transaction)

	return nil
}

// only used in tests

func (m *MockImpl) Reset() {
	m.recording = make([]Transaction, 0)
	m.simulateGetError = nil
	m.simulateAddError = nil
}

func (m *MockImpl) Recording() []Transaction {
	return m.recording
}

func (m *MockImpl) SimulateAddError(err error) {
	m.simulateAddError = err
}

func (m *MockImpl) InjectTransaction(_ context.Context, transaction Transaction) error {
	existingTransactions, ok := m.data[transaction.DebitorID]
	if !ok {
		existingTransactions = make([]Transaction, 0)
	}

	transactions := append(existingTransactions, transaction)
	m.data[transaction.DebitorID] = transactions

	return nil
}
