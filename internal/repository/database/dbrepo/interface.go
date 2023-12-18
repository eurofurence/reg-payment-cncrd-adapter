package dbrepo

import (
	"context"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/entity"
)

type Repository interface {
	Open() error
	Close()
	Migrate() error

	WriteProtocolEntry(ctx context.Context, e *entity.ProtocolEntry) error
}
