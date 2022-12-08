package paymentlinksrv

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/api/v1/cncrdapi"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/concardis"
)

func (i *Impl) HandleWebhook(ctx context.Context, webhook cncrdapi.WebhookEventDto) error {
	aulogging.Logger.Ctx(ctx).Info().Printf("webhook id=%d invoice.paymentRequestId=%d invoice.referenceId=%s", webhook.Transaction.Id, webhook.Transaction.Invoice.PaymentRequestId, webhook.Transaction.Invoice.ReferenceId)

	paylinkId, err := idValidate(webhook.Transaction.Invoice.PaymentRequestId)
	if err != nil {
		return err
	}

	paylink, err := concardis.Get().QueryPaymentLink(ctx, paylinkId)
	if err != nil {
		return err
	}

	aulogging.Logger.Ctx(ctx).Info().Printf("paylink id=%d ref=%s status=%s", paylink.ID, paylink.ReferenceID, paylink.Status)

	// TODO verify that the requestId in the webhook matches the retrieved paylink

	// TODO actually handle the transaction(s) in the paylink response (update payment to status pending in payment service, or failing that create one?)

	return nil
}

func idValidate(value int64) (uint, error) {
	if value < 1 {
		return 0, WebhookValidationErr
	}
	return uint(value), nil
}
