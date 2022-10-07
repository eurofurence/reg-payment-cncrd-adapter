package paylinkctl

import (
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/web/util/ctlutil"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/web/util/ctxvalues"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Create(server chi.Router) {
	server.Post("/api/rest/v1/paylinks", createPaylinkHandler)
}

func createPaylinkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if !ctxvalues.HasApiToken(ctx) {
		ctlutil.UnauthenticatedError(ctx, w, r, "you must be logged in for this operation", "anonymous access attempt")
		return
	}

}
