package paymentlinksrv

import (
	"context"
	"errors"
	"net/url"

	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/api/v1/cncrdapi"
)

type PaymentLinkService interface {
	// ValidatePaymentLinkRequest checks the cncrdapi.PaymentLinkRequestDto for validity.
	//
	// The returned url.Values contains detailed error messages that can be used to construct a meaningful response.
	// It is nil if no validation errors were encountered. Any errors encountered are also logged.
	ValidatePaymentLinkRequest(ctx context.Context, data cncrdapi.PaymentLinkRequestDto) url.Values

	// CreatePaymentLink expects an already validated cncrdapi.PaymentLinkRequestDto, and makes a downstream
	// request to create a payment link, returning the cncrdapi.PaymentLinkDto with all its information and the
	// id under which to manage the payment link.
	CreatePaymentLink(ctx context.Context, request cncrdapi.PaymentLinkRequestDto) (cncrdapi.PaymentLinkDto, uint, error)

	// GetPaymentLink obtains the payment link information from the downstream api.
	GetPaymentLink(ctx context.Context, id uint) (cncrdapi.PaymentLinkDto, error)

	// DeletePaymentLink asks the downstream api to delete the given payment link.
	DeletePaymentLink(ctx context.Context, id uint) error

	// HandleWebhook requests the payment link referenced in the webhook data and reacts to any new payments
	HandleWebhook(ctx context.Context, webhook cncrdapi.WebhookEventDto) error

	// SendErrorNotifyMail notifies us about unexpected conditions in this service so we can look at the logs
	SendErrorNotifyMail(ctx context.Context, operation string, referenceId string, status string) error
}

var (
	WebhookValidationErr    = errors.New("webhook referenced invalid invoice id, must be positive integer")
	WebhookRefIdMismatchErr = errors.New("webhook reference_id differes from paylink reference_id")
)
