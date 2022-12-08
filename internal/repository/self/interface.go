package self

import (
	"context"
	"errors"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/api/v1/cncrdapi"
)

type Self interface {
	CallWebhook(ctx context.Context, event cncrdapi.WebhookEventDto) error
}

var (
	DownstreamError = errors.New("downstream unavailable - see log for details")
)
