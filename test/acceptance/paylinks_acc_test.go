package acceptance

import (
	"github.com/eurofurence/reg-payment-cncrd-adapter/docs"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestCreatePaylink_Anonymous(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given an unauthenticated user")
	token := tstNoToken()

	docs.When("when they attempt to create a paylink")
	requestBody := tstBuildValidPaymentLinkRequest("pc1")
	response := tstPerformPost("/api/rest/v1/paylinks", tstRenderJson(requestBody), token)

	docs.Then("then the request is denied as unauthenticated (401) with the appropriate error message")
	tstRequireErrorResponse(t, response, http.StatusUnauthorized, "auth.unauthorized", "you must be logged in for this operation")

	docs.Then("and no requests to the payment provider have been made")
	require.Empty(t, concardisMock.Recording())
}
