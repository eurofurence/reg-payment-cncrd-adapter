package concardis

import (
	"context"
	"errors"
	"fmt"
	aurestbreaker "github.com/StephanHCB/go-autumn-restclient-circuitbreaker/implementation/breaker"
	aurestclientapi "github.com/StephanHCB/go-autumn-restclient/api"
	auresthttpclient "github.com/StephanHCB/go-autumn-restclient/implementation/httpclient"
	aurestlogging "github.com/StephanHCB/go-autumn-restclient/implementation/requestlogging"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/config"
	"net/http"
	"net/url"
	"time"
)

type Impl struct {
	client       aurestclientapi.Client
	baseUrl      string
	instanceName string
}

func newClient() (ConcardisDownstream, error) {
	httpClient, err := auresthttpclient.New(0, nil, nil)
	if err != nil {
		return nil, err
	}

	requestLoggingClient := aurestlogging.New(httpClient)

	circuitBreakerClient := aurestbreaker.New(requestLoggingClient,
		"concardis-downstream-breaker",
		10,
		2*time.Minute,
		30*time.Second,
		15*time.Second,
	)

	return &Impl{
		client:       circuitBreakerClient,
		baseUrl:      config.ConcardisDownstreamBaseUrl(),
		instanceName: config.ConcardisInstanceName(),
	}, nil
}

type createLowlevelResponseBody struct {
	Status string               `json:"status"`
	Data   []PaymentLinkCreated `json:"data"`
}

type createLowlevelRequestBody struct {
	PaymentLinkCreateRequest

	ApiSignature string `json:"ApiSignature"`
}

func (i *Impl) CreatePaymentLink(ctx context.Context, request PaymentLinkCreateRequest) (PaymentLinkCreated, error) {
	requestUrl := fmt.Sprintf("%s/v1.0/Invoice/?instance=%s", i.baseUrl, url.QueryEscape(i.instanceName))
	requestDto := createLowlevelRequestBody{
		PaymentLinkCreateRequest: request,
		ApiSignature:             "TODO",
	}
	bodyDto := createLowlevelResponseBody{}
	response := aurestclientapi.ParsedResponse{
		Body: &bodyDto,
	}
	if err := i.client.Perform(ctx, http.MethodPost, requestUrl, requestDto, &response); err != nil {
		return PaymentLinkCreated{}, err
	}
	if response.Status >= 300 {
		return PaymentLinkCreated{}, errors.New("TODO some sensible error")
	}
	if bodyDto.Status != "success" {
		return PaymentLinkCreated{}, errors.New("TODO some sensible error")
	}
	if len(bodyDto.Data) != 1 {
		return PaymentLinkCreated{}, errors.New("TODO some sensible error")
	}
	return bodyDto.Data[0], nil
}

type queryLowlevelResponseBody struct {
	Status string                     `json:"status"`
	Data   []PaymentLinkQueryResponse `json:"data"`
}

type queryLowlevelRequestBody struct {
	ApiSignature string `json:"ApiSignature"`
}

func (i *Impl) QueryPaymentLink(ctx context.Context, id uint) (PaymentLinkQueryResponse, error) {
	requestUrl := fmt.Sprintf("%s/v1.0/Invoice/?instance=%s", i.baseUrl, url.QueryEscape(i.instanceName))
	requestDto := queryLowlevelRequestBody{
		ApiSignature: "TODO",
	}
	bodyDto := queryLowlevelResponseBody{}
	response := aurestclientapi.ParsedResponse{
		Body: &bodyDto,
	}
	if err := i.client.Perform(ctx, http.MethodGet, requestUrl, requestDto, &response); err != nil {
		return PaymentLinkQueryResponse{}, err
	}
	if response.Status >= 300 {
		return PaymentLinkQueryResponse{}, errors.New("TODO some sensible error")
	}
	if bodyDto.Status != "success" {
		return PaymentLinkQueryResponse{}, errors.New("TODO some sensible error")
	}
	if len(bodyDto.Data) != 1 {
		return PaymentLinkQueryResponse{}, errors.New("TODO some sensible error")
	}
	return bodyDto.Data[0], nil
}

type deleteLowlevelResponseBody struct {
	Status string `json:"status"`
}

type deleteLowlevelRequestBody struct {
	ApiSignature string `json:"ApiSignature"`
}

func (i *Impl) DeletePaymentLink(ctx context.Context, id uint) error {
	requestUrl := fmt.Sprintf("%s/v1.0/Invoice/?instance=%s", i.baseUrl, url.QueryEscape(i.instanceName))
	requestDto := deleteLowlevelRequestBody{
		ApiSignature: "TODO",
	}
	bodyDto := deleteLowlevelResponseBody{}
	response := aurestclientapi.ParsedResponse{
		Body: &bodyDto,
	}
	if err := i.client.Perform(ctx, http.MethodGet, requestUrl, requestDto, &response); err != nil {
		return err
	}
	if response.Status >= 300 {
		return errors.New("TODO some sensible error")
	}
	if bodyDto.Status != "success" {
		return errors.New("TODO some sensible error")
	}
	return nil
}

func (i *Impl) QueryTransactions(ctx context.Context, timeGreaterThan time.Time, timeLessThan time.Time) ([]TransactionData, error) {
	//TODO implement me
	panic("implement me")
}
