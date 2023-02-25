package paymentlinksrv

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/api/v1/cncrdapi"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/concardis"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/paymentservice"
)

const isoDateFormat = "2006-01-02"

func (i *Impl) HandleWebhook(ctx context.Context, webhook cncrdapi.WebhookEventDto) error {
	aulogging.Logger.Ctx(ctx).Info().Printf("webhook id=%d invoice.paymentRequestId=%d invoice.referenceId=%s", webhook.Transaction.Id, webhook.Transaction.Invoice.PaymentRequestId, webhook.Transaction.Invoice.ReferenceId)

	paylinkId, err := idValidate(webhook.Transaction.Invoice.PaymentRequestId)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Error().Printf("webhook called with invalid paylink ID. id=%d", webhook.Transaction.Invoice.PaymentRequestId)
		return err
	}

	paylink, err := concardis.Get().QueryPaymentLink(ctx, paylinkId)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Error().Printf("can't query payment link from concardis. err=%s", err.Error())
		return err
	}

	if paylink.ReferenceID != webhook.Transaction.Invoice.ReferenceId {
		// webhook data claimed it was about ref_id A, but the paylink is for ref_id B
		aulogging.Logger.Ctx(ctx).Error().Printf("webhook reference_id mismatch, ref_id in webhook=%s, ref_id in paylink data=%s", webhook.Transaction.Invoice.ReferenceId, paylink.ReferenceID)
		return WebhookRefIdMismatchErr
	}

	aulogging.Logger.Ctx(ctx).Info().Printf("webhook call for paylink id=%d ref=%s status=%s amount=%d", paylink.ID, paylink.ReferenceID, paylink.Status, paylink.Amount)

	// fetch transaction data from payment service
	transaction, err := paymentservice.Get().GetTransactionByReferenceId(ctx, paylink.ReferenceID)
	if err != nil {
		if err == paymentservice.NotFoundError {
			// transaction not found in the payment service -> create one.
			// Note: this should never happen, but we try to recover because someone paid us money for somthing.
			aulogging.Logger.Ctx(ctx).Error().Printf("webhook reference_id not found in payment service. Creating new transaction. reference_id=%s", paylink.ReferenceID)

			return createTransaction(ctx, paylink)
		} else {
			aulogging.Logger.Ctx(ctx).Error().Printf("error fetching transaction from payment service. err=%s", err.Error())
			return err
		}
	}

	// matching transaction was found in the payment service database.
	// update the values with data from Concardis.
	return updateTransaction(ctx, paylink, transaction)
}

func createTransaction(ctx context.Context, paylink concardis.PaymentLinkQueryResponse) error {
	debitor_id, err := debitorIdFromReferenceID(paylink.ReferenceID)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().Printf("webhook couldn't parse debitor_id from reference_id. reference_id=%s", paylink.ReferenceID)
		// we log a warning, but we continue anyway
	}

	today := time.Now().Format(isoDateFormat)

	transaction := paymentservice.Transaction{
		ID:        paylink.ReferenceID,
		DebitorID: debitor_id,
		Type:      paymentservice.Payment,
		Method:    paymentservice.Credit, // we use paylink for credit cards only, atm.
		Amount: paymentservice.Amount{
			GrossCent: paylink.Amount,
			Currency:  paylink.Currency,
			VatRate:   0, // TODO should set from payload
		},
		Comment:       "Auto-created by cncrd adapter because the reference_id could not be found in the payment service.",
		Status:        paymentservice.Pending,
		EffectiveDate: today, // TODO: this should be in the payload
		DueDate:       today,
		// omitting Deletion
	}

	err = paymentservice.Get().AddTransaction(ctx, transaction)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Error().Printf("webhook could not create transaction in payment service! (we don't know why we received this money, and we couldn't add the transaction to the database either!) reference_id=%s", paylink.ReferenceID)
	}
	return err
}

func updateTransaction(ctx context.Context, paylink concardis.PaymentLinkQueryResponse, transaction paymentservice.Transaction) error {
	transaction.Amount.GrossCent = paylink.Amount
	transaction.Amount.Currency = paylink.Currency
	transaction.Status = paymentservice.Pending           // TODO fail if already valid and values do not match (admin might have done this in the mean time)
	transaction.EffectiveDate = transaction.EffectiveDate // TODO: this should be in the payload

	err := paymentservice.Get().UpdateTransaction(ctx, transaction)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Error().Printf("webhook unable to update upstream transaction. reference_id=%s", paylink.ReferenceID)
		return err
	}

	return nil
}

func debitorIdFromReferenceID(ref_id string) (uint, error) {
	// reference_id is generated internally in the payment service.
	// See  reg-payment-service/internal/interaction/transaction.go:generateTransactionID()

	// The format is:  "%s-%06d-%s-%s"
	// Fields:
	//   - prefix (hopefully without hyphens)
	//   - debitor_id
	//   - timestamp in format "0102-150405" (hyphen!)
	//   - random digits

	tokens := strings.Split(ref_id, "-")
	if len(tokens) != 5 {
		return 0, errors.New("error parsing reference_id")
	}

	debitor_id, err := strconv.ParseUint(tokens[1], 10, 32)
	if err != nil {
		return 0, errors.New("error parsing debitor_id as uint: " + err.Error())
	}

	return uint(debitor_id), nil
}

func idValidate(value int64) (uint, error) {
	if value < 1 {
		return 0, WebhookValidationErr
	}
	return uint(value), nil
}
