package paymentlinksrv

import (
	"context"
	"net/url"

	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/api/v1/cncrdapi"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/concardis"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/config"
)

func (i *Impl) ValidatePaymentLinkRequest(ctx context.Context, data cncrdapi.PaymentLinkRequestDto) url.Values {
	errs := url.Values{}

	if data.DebitorId == 0 {
		errs.Add("debitor_id", "field must be a positive integer (the badge number to bill for)")
	}
	if data.AmountDue <= 0 {
		errs.Add("amount_due", "must be a positive integer (the amount to bill)")
	}
	if data.Currency != "EUR" {
		errs.Add("currency", "right now, only EUR is supported")
	}
	if data.VatRate < 0.0 || data.VatRate > 50.0 {
		errs.Add("vat_rate", "vat rate should be provided in percent and must be between 0.0 and 50.0")
	}

	if len(errs) == 0 {
		return nil
	} else {
		return errs
	}
}

func (i *Impl) CreatePaymentLink(ctx context.Context, data cncrdapi.PaymentLinkRequestDto) (cncrdapi.PaymentLinkDto, uint, error) {
	concardisRequest := i.concardisCreateRequestFromApiRequest(data)
	concardisResponse, err := concardis.Get().CreatePaymentLink(ctx, concardisRequest)
	if err != nil {
		return cncrdapi.PaymentLinkDto{}, 0, err
	}
	output := i.apiResponseFromConcardisResponse(concardisResponse, concardisRequest)
	return output, concardisResponse.ID, nil
}

func (i *Impl) concardisCreateRequestFromApiRequest(data cncrdapi.PaymentLinkRequestDto) concardis.PaymentLinkCreateRequest {
	return concardis.PaymentLinkCreateRequest{
		Title:       config.InvoiceTitle(),
		Description: config.InvoiceDescription(),
		PSP:         1,
		ReferenceId: data.ReferenceId,
		OrderId:     data.ReferenceId,
		Purpose:     config.InvoicePurpose(),
		Amount:      data.AmountDue,
		VatRate:     data.VatRate,
		Currency:    data.Currency,
		SKU:         "registration",
	}
}

func (i *Impl) apiResponseFromConcardisResponse(response concardis.PaymentLinkCreated, request concardis.PaymentLinkCreateRequest) cncrdapi.PaymentLinkDto {
	return cncrdapi.PaymentLinkDto{
		Title:       request.Title,
		Description: request.Description,
		ReferenceId: response.ReferenceID,
		Purpose:     request.Purpose,
		AmountDue:   request.Amount,
		AmountPaid:  0,
		Currency:    request.Currency,
		VatRate:     request.VatRate,
		Link:        response.Link,
	}
}

func (i *Impl) GetPaymentLink(ctx context.Context, id uint) (cncrdapi.PaymentLinkDto, error) {
	data, err := concardis.Get().QueryPaymentLink(ctx, id)
	if err != nil {
		return cncrdapi.PaymentLinkDto{}, err
	}

	// TODO lots of missing fields, can we get them from downstream?

	result := cncrdapi.PaymentLinkDto{
		ReferenceId: data.ReferenceID,
		Purpose:     data.Purpose["1"],
		AmountDue:   data.Amount,
		AmountPaid:  0,
		Currency:    data.Currency,
		Link:        data.Link,
	}

	return result, nil
}

func (i *Impl) DeletePaymentLink(ctx context.Context, id uint) error {
	err := concardis.Get().DeletePaymentLink(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
