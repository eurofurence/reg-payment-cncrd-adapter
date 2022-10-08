package paymentlinksrv

import (
	"context"
	"fmt"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/api/v1/cncrdapi"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/concardis"
)

func (i *Impl) CreatePaymentLink(ctx context.Context, data cncrdapi.PaymentLinkRequestDto) (cncrdapi.PaymentLinkDto, int64, error) {
	// TODO validation
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
		// TODO implement with help of configuration
		Title:       "some page title",
		Description: "some page description",
		PSP:         1,
		ReferenceId: fmt.Sprintf("144823ad-%06d", data.DebitorId),
		Purpose:     "some payment purpose",
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
