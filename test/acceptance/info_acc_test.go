package acceptance

import (
	"github.com/eurofurence/reg-payment-cncrd-adapter/docs"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/web/util/media"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given an unauthenticated user")

	docs.When("when they perform GET on the health endpoint")
	response := tstPerformGet("/info/health", tstNoToken())

	docs.Then("then OK is returned, and no further information is available")
	require.Equal(t, http.StatusOK, response.status, "unexpected http response status")
	require.Equal(t, media.ContentTypeTextPlain, response.contentType, "unexpected response content type")
	require.Equal(t, "OK", response.body, "unexpected response from health endpoint")
}

func TestErrorFallback(t *testing.T) {
	tstSetup(tstConfigFile)
	defer tstShutdown()

	docs.Given("given an unauthenticated user")

	docs.When("when they perform GET on an unimplemented endpoint")
	response := tstPerformGet("/info/does-not-exist", tstNoToken())

	docs.Then("then they receive a 404 error")
	require.Equal(t, http.StatusNotFound, response.status, "unexpected http response status")
}
