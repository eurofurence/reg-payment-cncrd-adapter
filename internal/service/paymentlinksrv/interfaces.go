package paymentlinksrv

import (
	"context"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/api/v1/cncrdapi"
)

type PaymentLinkService interface {
	CreatePaymentLink(ctx context.Context, request cncrdapi.PaymentLinkRequestDto) (cncrdapi.PaymentLinkDto, int64, error)
}
