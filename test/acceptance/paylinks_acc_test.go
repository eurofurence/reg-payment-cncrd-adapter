package acceptance

import (
	"github.com/eurofurence/reg-payment-cncrd-adapter/docs"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestCreatePaylink_Success(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given a caller who supplies a correct api token")
	token := tstValidApiToken()

	docs.When("when they attempt to create a payment link with valid information")
	requestBody := tstBuildValidPaymentLinkRequest()
	response := tstPerformPost("/api/rest/v1/paylinks", tstRenderJson(requestBody), token)

	docs.Then("then the request is successful and the response is as expected")
	tstRequirePaymentLinkResponse(t, response, http.StatusCreated)

	docs.Then("and the expected request for a payment link has been made")
	tstRequireConcardisRecording(t,
		"CreatePaymentLink {some page title some page description 1 144823ad-000001 some payment purpose 390 19 EUR registration}",
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
	tstRequireErrorResponse(t, response, http.StatusBadRequest, "body.parse.error", nil)

	docs.Then("and no requests to the payment provider have been made")
	require.Empty(t, concardisMock.Recording())
}

// TODO business validation stuff

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
