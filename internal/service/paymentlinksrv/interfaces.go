package paymentlinksrv

import (
	"context"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/api/v1/cncrdapi"
	"net/url"
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
	CreatePaymentLink(ctx context.Context, request cncrdapi.PaymentLinkRequestDto) (cncrdapi.PaymentLinkDto, int64, error)

	// GetPaymentLink obtains the payment link information from the downstream api.
	GetPaymentLink(ctx context.Context, id uint) (cncrdapi.PaymentLinkDto, error)

	// DeletePaymentLink asks the downstream api to delete the given payment link.
	DeletePaymentLink(ctx context.Context, id uint) error
}
