// Package simulatorctl implements a local simulator for the api.
//
// If you do not set up a url to the payment provider in the configuration, this will be what the service
// talks to instead. Generated simulated pay links will be handled by this implementation as well, and
// lead to the expected webhook call.
package simulatorctl

import (
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/service/paymentlinksrv"
	"github.com/go-chi/chi/v5"
	"net/http"
)

var paymentLinkService paymentlinksrv.PaymentLinkService

func Create(server chi.Router, paymentLinkSrv paymentlinksrv.PaymentLinkService) {
	paymentLinkService = paymentLinkSrv
	server.Get("/simulator/{referenceId}", useSimulator)
}

func useSimulator(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
}
