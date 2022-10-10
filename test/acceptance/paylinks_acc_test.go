package acceptance

import (
	"github.com/eurofurence/reg-payment-cncrd-adapter/docs"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/api/v1/cncrdapi"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/concardis"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

// --- create ---

func TestCreatePaylink_Success(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given a caller who supplies a correct api token")
	token := tstValidApiToken()

	docs.When("when they attempt to create a payment link with valid information")
	requestBody := tstBuildValidPaymentLinkRequest()
	response := tstPerformPost("/api/rest/v1/paylinks", tstRenderJson(requestBody), token)

	docs.Then("then the request is successful and the response is as expected")
	tstRequirePaymentLinkResponse(t, response, http.StatusCreated, tstBuildValidPaymentLink())

	docs.Then("and the expected request for a payment link has been made")
	tstRequireConcardisRecording(t,
		"CreatePaymentLink {some page title some page description 1 221216-122218-000001 some payment purpose 390 19 EUR registration}",
	)
}

func TestCreatePaylink_InvalidJson(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given a caller who supplies a correct api token")
	token := tstValidApiToken()

	docs.When("when they attempt to create a payment link but supply an invalid json body")
	response := tstPerformPost("/api/rest/v1/paylinks", "{{{::}{}{}", token)

	docs.Then("then the request is denied with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "paylink.parse.error", nil)

	docs.Then("and no requests to the payment provider have been made")
	require.Empty(t, concardisMock.Recording())
}

func TestCreatePaylink_ValidJsonWrongFields(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given a caller who supplies a correct api token")
	token := tstValidApiToken()

	docs.When("when they attempt to create a payment link but supply a json body with the wrong fields")
	response := tstPerformPost("/api/rest/v1/paylinks", `{"hello":"kitty"}`, token)

	docs.Then("then the request is denied with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "paylink.parse.error", nil)

	docs.Then("and no requests to the payment provider have been made")
	require.Empty(t, concardisMock.Recording())
}

func TestCreatePaylink_InvalidData(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given a caller who supplies a correct api token")
	token := tstValidApiToken()

	docs.When("when they attempt to create a payment link but supply invalid field values")
	requestBody := cncrdapi.PaymentLinkRequestDto{
		DebitorId: 0,
		AmountDue: -53,
		Currency:  "CHF",
		VatRate:   -33.3,
	}
	response := tstPerformPost("/api/rest/v1/paylinks", tstRenderJson(requestBody), token)

	docs.Then("then the request is denied with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "paylink.data.invalid", url.Values{
		"amount_due": []string{"must be a positive integer (the amount to bill)"},
		"currency":   []string{"right now, only EUR is supported"},
		"debitor_id": []string{"field must be a positive integer (the badge number to bill for)"},
		"vat_rate":   []string{"vat rate should be provided in percent and must be between 0.0 and 50.0"},
	})

	docs.Then("and no requests to the payment provider have been made")
	require.Empty(t, concardisMock.Recording())
}

func TestCreatePaylink_Anonymous(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given an unauthenticated caller")
	token := tstNoToken()

	docs.When("when they attempt to create a payment link")
	requestBody := tstBuildValidPaymentLinkRequest()
	response := tstPerformPost("/api/rest/v1/paylinks", tstRenderJson(requestBody), token)

	docs.Then("then the request is denied as unauthenticated (401) with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")

	docs.Then("and no requests to the payment provider have been made")
	require.Empty(t, concardisMock.Recording())
}

func TestCreatePaylink_WrongToken(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given a caller who supplies a wrong api token")
	token := tstInvalidApiToken()

	docs.When("when they attempt to create a payment link")
	requestBody := tstBuildValidPaymentLinkRequest()
	response := tstPerformPost("/api/rest/v1/paylinks", tstRenderJson(requestBody), token)

	docs.Then("then the request is denied as unauthenticated (401) with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "invalid api token")

	docs.Then("and no requests to the payment provider have been made")
	require.Empty(t, concardisMock.Recording())
}

func TestCreatePaylink_DownstreamError(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given a caller who supplies a correct api token")
	token := tstValidApiToken()

	docs.When("when they attempt to create a payment link with valid information while the downstream api is down")
	concardisMock.SimulateError(concardis.DownstreamError)
	requestBody := tstBuildValidPaymentLinkRequest()
	response := tstPerformPost("/api/rest/v1/paylinks", tstRenderJson(requestBody), token)

	docs.Then("then the request fails with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusBadGateway, "paylink.downstream.error", nil)
}

// --- get ---

func TestGetPaylink_Success(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given a caller who supplies a correct api token")
	token := tstValidApiToken()

	docs.When("when they attempt to get an existing payment link")
	response := tstPerformGet("/api/rest/v1/paylinks/42", token)

	docs.Then("then the request is successful and the response is as expected")
	tstRequirePaymentLinkResponse(t, response, http.StatusOK, tstBuildValidPaymentLinkGetResponse())

	docs.Then("and the expected request for a payment link has been made")
	tstRequireConcardisRecording(t,
		"QueryPaymentLink 42",
	)
}

func TestGetPaylink_InvalidId(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given a caller who supplies a correct api token")
	token := tstValidApiToken()

	docs.When("when they attempt to get a payment link but supply an invalid id")
	response := tstPerformGet("/api/rest/v1/paylinks/%2f%4c", token)

	docs.Then("then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "paylink.id.invalid", nil)

	docs.Then("and no requests to the payment provider have been made")
	require.Empty(t, concardisMock.Recording())
}

func TestGetPaylink_NotFound(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given a caller who supplies a correct api token")
	token := tstValidApiToken()

	docs.When("when they attempt to get a payment link but supply an id that does not exist")
	response := tstPerformGet("/api/rest/v1/paylinks/13", token)

	docs.Then("then the request fails with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusNotFound, "paylink.id.notfound", nil)

	docs.Then("and the expected request for a payment link has been made")
	tstRequireConcardisRecording(t,
		"QueryPaymentLink 13",
	)
}

func TestGetPaylink_Anonymous(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given an unauthenticated caller")
	token := tstNoToken()

	docs.When("when they attempt to get a payment link")
	response := tstPerformGet("/api/rest/v1/paylinks/42", token)

	docs.Then("then the request is denied as unauthenticated (401) with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")

	docs.Then("and no requests to the payment provider have been made")
	require.Empty(t, concardisMock.Recording())
}

func TestGetPaylink_WrongToken(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given a caller who supplies a wrong api token")
	token := tstInvalidApiToken()

	docs.When("when they attempt to get a payment link")
	response := tstPerformGet("/api/rest/v1/paylinks/42", token)

	docs.Then("then the request is denied as unauthenticated (401) with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "invalid api token")

	docs.Then("and no requests to the payment provider have been made")
	require.Empty(t, concardisMock.Recording())
}

func TestGetPaylink_DownstreamError(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given a caller who supplies a correct api token")
	token := tstValidApiToken()

	docs.When("when they attempt to get a payment link while the downstream api is down")
	concardisMock.SimulateError(concardis.DownstreamError)
	response := tstPerformGet("/api/rest/v1/paylinks/42", token)

	docs.Then("then the request fails with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusBadGateway, "paylink.downstream.error", nil)
}
