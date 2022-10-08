package paylinkctl

import (
	"context"
	"encoding/json"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/api/v1/cncrdapi"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/service/paymentlinksrv"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/web/util/ctlutil"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/web/util/ctxvalues"
	"github.com/go-chi/chi/v5"
	"github.com/go-http-utils/headers"
	"net/http"
)

var paymentLinkService paymentlinksrv.PaymentLinkService

func init() {
	paymentLinkService = &paymentlinksrv.Impl{}
}

func Create(server chi.Router) {
	server.Post("/api/rest/v1/paylinks", createPaylinkHandler)
}

func createPaylinkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if !ctxvalues.HasApiToken(ctx) {
		ctlutil.UnauthenticatedError(ctx, w, r, "you must be logged in for this operation", "anonymous access attempt")
		return
	}

	request, err := parseBodyToPaymentLinkRequestDto(ctx, w, r)
	if err != nil {
		return
	}

	dto, id, err := paymentLinkService.CreatePaymentLink(ctx, request)
	if err != nil {
		// TODO handle specific errors like validation
		ctlutil.UnexpectedError(ctx, w, r, err)
	}

	w.Header().Set(headers.Location, fmt.Sprintf("/api/rest/v1/paylinks/%d", id))
	w.WriteHeader(http.StatusCreated)
	ctlutil.WriteJson(ctx, w, dto)
}

func parseBodyToPaymentLinkRequestDto(ctx context.Context, w http.ResponseWriter, r *http.Request) (cncrdapi.PaymentLinkRequestDto, error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	dto := cncrdapi.PaymentLinkRequestDto{}
	err := decoder.Decode(&dto)
	if err != nil {
		bodyParseErrorHandler(ctx, w, r, err)
	}
	return dto, err
}

func bodyParseErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("body could not be parsed: %s", err.Error())
	ctlutil.ErrorHandler(ctx, w, r, "body.parse.error", http.StatusBadRequest, nil)
}
