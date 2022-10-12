package paylinkctl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/api/v1/cncrdapi"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/concardis"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/service/paymentlinksrv"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/web/util/ctlutil"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/web/util/ctxvalues"
	"github.com/go-chi/chi/v5"
	"github.com/go-http-utils/headers"
	"net/http"
	"net/url"
	"strconv"
)

var paymentLinkService paymentlinksrv.PaymentLinkService

func Create(server chi.Router, paymentLinkSrv paymentlinksrv.PaymentLinkService) {
	paymentLinkService = paymentLinkSrv

	server.Post("/api/rest/v1/paylinks", createPaylinkHandler)
	server.Get("/api/rest/v1/paylinks/{id}", getPaylinkHandler)
	server.Delete("/api/rest/v1/paylinks/{id}", deletePaylinkHandler)
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

	errs := paymentLinkService.ValidatePaymentLinkRequest(ctx, request)
	if errs != nil {
		paylinkRequestInvalidErrorHandler(ctx, w, r, errs)
		return
	}

	dto, id, err := paymentLinkService.CreatePaymentLink(ctx, request)
	if err != nil {
		if errors.Is(err, concardis.DownstreamError) {
			downstreamErrorHandler(ctx, w, r, err)
		} else {
			ctlutil.UnexpectedError(ctx, w, r, err)
		}
		return
	}

	w.Header().Set(headers.Location, fmt.Sprintf("/api/rest/v1/paylinks/%d", id))
	w.WriteHeader(http.StatusCreated)
	ctlutil.WriteJson(ctx, w, dto)
}

func getPaylinkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if !ctxvalues.HasApiToken(ctx) {
		ctlutil.UnauthenticatedError(ctx, w, r, "you must be logged in for this operation", "anonymous access attempt")
		return
	}

	id, err := idFromVars(ctx, w, r)
	if err != nil {
		return
	}

	dto, err := paymentLinkService.GetPaymentLink(ctx, id)
	if err != nil {
		if errors.Is(err, concardis.DownstreamError) {
			downstreamErrorHandler(ctx, w, r, err)
		} else if errors.Is(err, concardis.NoSuchID404Error) {
			paylinkNotFoundErrorHandler(ctx, w, r, id)
		} else {
			ctlutil.UnexpectedError(ctx, w, r, err)
		}
		return
	}

	ctlutil.WriteJson(ctx, w, dto)
}

func deletePaylinkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if !ctxvalues.HasApiToken(ctx) {
		ctlutil.UnauthenticatedError(ctx, w, r, "you must be logged in for this operation", "anonymous access attempt")
		return
	}

	id, err := idFromVars(ctx, w, r)
	if err != nil {
		return
	}

	err = paymentLinkService.DeletePaymentLink(ctx, id)
	if err != nil {
		if errors.Is(err, concardis.DownstreamError) {
			downstreamErrorHandler(ctx, w, r, err)
		} else if errors.Is(err, concardis.NoSuchID404Error) {
			paylinkNotFoundErrorHandler(ctx, w, r, id)
		} else {
			ctlutil.UnexpectedError(ctx, w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseBodyToPaymentLinkRequestDto(ctx context.Context, w http.ResponseWriter, r *http.Request) (cncrdapi.PaymentLinkRequestDto, error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	dto := cncrdapi.PaymentLinkRequestDto{}
	err := decoder.Decode(&dto)
	if err != nil {
		paylinkRequestParseErrorHandler(ctx, w, r, err)
	}
	return dto, err
}

func idFromVars(ctx context.Context, w http.ResponseWriter, r *http.Request) (uint, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		invalidPaylinkIdErrorHandler(ctx, w, r, idStr)
	}
	return uint(id), err
}

func paylinkRequestParseErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("paylink body could not be parsed: %s", err.Error())
	ctlutil.ErrorHandler(ctx, w, r, "paylink.parse.error", http.StatusBadRequest, nil)
}

func paylinkRequestInvalidErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, validationErrors url.Values) {
	// validation already logged each individual error
	ctlutil.ErrorHandler(ctx, w, r, "paylink.data.invalid", http.StatusBadRequest, validationErrors)
}

func downstreamErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("downstream error: %s", err.Error())
	ctlutil.ErrorHandler(ctx, w, r, "paylink.downstream.error", http.StatusBadGateway, nil)
}

func invalidPaylinkIdErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, id string) {
	aulogging.Logger.Ctx(ctx).Warn().Printf("received invalid paylink id '%s'", url.QueryEscape(id))
	ctlutil.ErrorHandler(ctx, w, r, "paylink.id.invalid", http.StatusBadRequest, url.Values{})
}

func paylinkNotFoundErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, id uint) {
	aulogging.Logger.Ctx(ctx).Warn().Printf("paylink id %d not found", id)
	ctlutil.ErrorHandler(ctx, w, r, "paylink.id.notfound", http.StatusNotFound, url.Values{})
}
