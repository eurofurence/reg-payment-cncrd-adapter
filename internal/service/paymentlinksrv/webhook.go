package paymentlinksrv

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/api/v1/cncrdapi"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/concardis"
	"strconv"
)

func (i *Impl) HandleWebhook(ctx context.Context, webhook cncrdapi.WebhookEventDto) error {
	aulogging.Logger.Ctx(ctx).Info().Printf("webhook id=%d invoice.number=%s", webhook.Transaction.Id, webhook.Transaction.Invoice.Number)

	paylinkId, err := idFromStr(webhook.Transaction.Invoice.Number)
	if err != nil {
		return err
	}

	paylink, err := concardis.Get().QueryPaymentLink(ctx, paylinkId)
	if err != nil {
		return err
	}

	aulogging.Logger.Ctx(ctx).Info().Printf("paylink id=%d ref=%s status=%s", paylink.ID, paylink.ReferenceID, paylink.Status)

	// TODO actually handle payments

	return nil
}

func idFromStr(value string) (uint, error) {
	idInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	if idInt < 1 {
		return 0, WebhookValidationErr
	}
	return uint(idInt), nil
}
