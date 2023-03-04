package acceptance

import (
	"github.com/eurofurence/reg-payment-cncrd-adapter/docs"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/concardis"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/paymentservice"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestWebhook_Success_TolerantReader(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given an anonymous caller who knows the secret url")
	url := "/api/rest/v1/webhook/demosecret"

	docs.When("when they trigger our webhook endpoint with valid information with lots of extra fields which we ignore")
	response := tstPerformPost(url, tstBuildValidWebhookRequest(), tstNoToken())

	docs.Then("then the request is successful")
	require.Equal(t, http.StatusOK, response.status)

	docs.Then("and the expected downstream requests have been made")
	tstRequireConcardisRecording(t,
		"QueryPaymentLink 42",
	)

	docs.Then("and the expected requests to the payment service have been made")
	tstRequirePaymentServiceRecording(t, []paymentservice.Transaction{
		{
			ID: "mock-transaction-id",
			Amount: paymentservice.Amount{
				Currency:  "EUR",
				GrossCent: 390,
			},
			Status: "pending",
		},
	})
}

func TestWebhook_Success_Status(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given an anonymous caller who knows the secret url")
	url := "/api/rest/v1/webhook/demosecret"

	docs.When("when they trigger our webhook endpoint with valid information with lots of extra fields which we ignore")
	response := tstPerformPost(url, tstBuildValidWebhookRequest(), tstNoToken())

	docs.Then("then the request is successful")
	require.Equal(t, http.StatusOK, response.status)

	docs.Then("and the expected downstream requests have been made")
	tstRequireConcardisRecording(t,
		"QueryPaymentLink 42",
	)

	docs.Then("and the expected requests to the payment service have been made")
	tstRequirePaymentServiceRecording(t, []paymentservice.Transaction{
		{
			ID: "mock-transaction-id",
			Amount: paymentservice.Amount{
				Currency:  "EUR",
				GrossCent: 390,
			},
			Status: "pending",
		},
	})
}

func TestWebhook_InvalidJson(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given an anonymous caller who knows the secret url")
	url := "/api/rest/v1/webhook/demosecret"

	docs.When("when they attempt to trigger our webhook endpoint with an invalid json body")
	response := tstPerformPost(url, `{{{{}}`, tstNoToken())

	docs.Then("then the request fails with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "webhook.parse.error", nil)
}

func TestWebhook_WrongSecret(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given an anonymous caller who does not know the secret url")
	url := "/api/rest/v1/webhook/wrongsecret"

	docs.When("when they attempt to trigger our webhook endpoint")
	response := tstPerformPost(url, tstBuildValidWebhookRequest(), tstNoToken())

	docs.Then("then the request fails with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", nil)
}

func TestWebhook_DownstreamError(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given an anonymous caller who knows the secret url")
	url := "/api/rest/v1/webhook/demosecret"

	docs.When("when they attempt to trigger our webhook endpoint while the downstream api is down")
	concardisMock.SimulateError(concardis.DownstreamError)
	response := tstPerformPost(url, tstBuildValidWebhookRequest(), tstNoToken())

	docs.Then("then the request fails with the appropriate error")
	tstRequireErrorResponse(t, response, http.StatusBadGateway, "webhook.downstream.error", nil)
}
