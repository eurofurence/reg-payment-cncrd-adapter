package concardis

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	aurestbreaker "github.com/StephanHCB/go-autumn-restclient-circuitbreaker/implementation/breaker"
	aurestclientapi "github.com/StephanHCB/go-autumn-restclient/api"
	auresthttpclient "github.com/StephanHCB/go-autumn-restclient/implementation/httpclient"
	aurestlogging "github.com/StephanHCB/go-autumn-restclient/implementation/requestlogging"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/config"
	"github.com/go-http-utils/headers"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Impl struct {
	client       aurestclientapi.Client
	baseUrl      string
	instanceName string
}

func requestManipulator(ctx context.Context, r *http.Request) {
	// even GET gets a body with the signature
	r.Header.Set(headers.ContentType, aurestclientapi.ContentTypeApplicationXWwwFormUrlencoded)
}

func newClient() (ConcardisDownstream, error) {
	httpClient, err := auresthttpclient.New(0, nil, requestManipulator)
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

// FixedSignatureValue set for automated contract testing only
var FixedSignatureValue string

const signatureKey = "ApiSignature"

func encode(key string, value string) string {
	return url.QueryEscape(key) + "=" + url.QueryEscape(value)
}

func signRequest(unsignedRequest string, instanceApiSecret string) string {
	// request parameters have to be in order for signature
	authenticator := hmac.New(sha256.New, []byte(instanceApiSecret))
	authenticator.Write([]byte(unsignedRequest))

	hashValue := authenticator.Sum([]byte{})
	signature := base64.StdEncoding.EncodeToString(hashValue)

	if FixedSignatureValue != "" {
		return FixedSignatureValue
	} else {
		return signature
	}
}

func buildCreateRequestBody(request PaymentLinkCreateRequest) string {
	var buf strings.Builder
	buf.WriteString(encode("title", request.Title) + "&")
	buf.WriteString(encode("description", request.Description) + "&")
	buf.WriteString(encode("psp", fmt.Sprintf("%d", request.PSP)) + "&")
	buf.WriteString(encode("referenceId", request.ReferenceId) + "&")
	buf.WriteString(encode("purpose", request.Purpose) + "&")
	buf.WriteString(encode("amount", fmt.Sprintf("%d", request.Amount)) + "&")
	buf.WriteString(encode("vatRate", fmt.Sprintf("%.1f", request.VatRate)) + "&")
	buf.WriteString(encode("currency", request.Currency) + "&")
	buf.WriteString(encode("sku", request.SKU) + "&")
	buf.WriteString(encode("preAuthorization", "0") + "&")
	buf.WriteString(encode("reservation", "0"))
	unsigned := buf.String()
	signature := signRequest(unsigned, config.ConcardisInstanceApiSecret())
	return unsigned + "&" + encode(signatureKey, signature)
}

func (i *Impl) CreatePaymentLink(ctx context.Context, request PaymentLinkCreateRequest) (PaymentLinkCreated, error) {
	requestUrl := fmt.Sprintf("%s/v1.0/Invoice/?instance=%s", i.baseUrl, url.QueryEscape(i.instanceName))
	requestBody := buildCreateRequestBody(request)
	bodyDto := createLowlevelResponseBody{}
	response := aurestclientapi.ParsedResponse{
		Body: &bodyDto,
	}
	if err := i.client.Perform(ctx, http.MethodPost, requestUrl, requestBody, &response); err != nil {
		return PaymentLinkCreated{}, err
	}
	if response.Status >= 300 {
		return PaymentLinkCreated{}, fmt.Errorf("unexpected response status %d", response.Status)
	}
	if bodyDto.Status != "success" {
		return PaymentLinkCreated{}, NotSuccessful
	}
	if len(bodyDto.Data) != 1 {
		return PaymentLinkCreated{}, fmt.Errorf("unexpected number of response body data array entries %d", len(bodyDto.Data))
	}
	return bodyDto.Data[0], nil
}

type queryLowlevelResponseBody struct {
	Status string                     `json:"status"`
	Data   []PaymentLinkQueryResponse `json:"data"`
}

func buildEmptyRequestBody() string {
	signature := signRequest("", config.ConcardisInstanceApiSecret())
	return encode(signatureKey, signature)
}

func (i *Impl) QueryPaymentLink(ctx context.Context, id uint) (PaymentLinkQueryResponse, error) {
	requestUrl := fmt.Sprintf("%s/v1.0/Invoice/%d/?instance=%s", i.baseUrl, id, url.QueryEscape(i.instanceName))
	requestBody := buildEmptyRequestBody()
	bodyDto := queryLowlevelResponseBody{}
	response := aurestclientapi.ParsedResponse{
		Body: &bodyDto,
	}
	if err := i.client.Perform(ctx, http.MethodGet, requestUrl, requestBody, &response); err != nil {
		return PaymentLinkQueryResponse{}, err
	}
	if response.Status >= 300 {
		return PaymentLinkQueryResponse{}, fmt.Errorf("unexpected response status %d", response.Status)
	}
	if bodyDto.Status != "success" {
		return PaymentLinkQueryResponse{}, NotSuccessful
	}
	if len(bodyDto.Data) != 1 {
		return PaymentLinkQueryResponse{}, fmt.Errorf("unexpected number of response body data array entries %d", len(bodyDto.Data))
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
	requestUrl := fmt.Sprintf("%s/v1.0/Invoice/%d/?instance=%s", i.baseUrl, id, url.QueryEscape(i.instanceName))
	requestBody := buildEmptyRequestBody()
	bodyDto := deleteLowlevelResponseBody{}
	response := aurestclientapi.ParsedResponse{
		Body: &bodyDto,
	}
	if err := i.client.Perform(ctx, http.MethodDelete, requestUrl, requestBody, &response); err != nil {
		return err
	}
	if response.Status >= 300 {
		return fmt.Errorf("unexpected response status %d", response.Status)
	}
	//if bodyDto.Status != "success" {
	//	return NotSuccessful
	//}
	return nil
}

func (i *Impl) QueryTransactions(ctx context.Context, timeGreaterThan time.Time, timeLessThan time.Time) ([]TransactionData, error) {
	//TODO implement me
	panic("implement me")
}
