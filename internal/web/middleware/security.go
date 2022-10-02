package middleware

import (
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/config"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/web/util/ctlutil"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/web/util/ctxvalues"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/web/util/media"
	"net/http"
)

// --- getting the values from the request ---

func fromApiTokenHeader(r *http.Request) string {
	return r.Header.Get(media.HeaderXApiKey)
}

// --- middleware validating the values and adding to context values ---

func TokenValidator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		apiTokenValue := fromApiTokenHeader(r)
		if apiTokenValue != "" {
			// ignore jwt if set (may still need to pass it through to other service)
			if apiTokenValue == config.FixedApiToken() {
				ctxvalues.SetApiToken(ctx, apiTokenValue)
				next.ServeHTTP(w, r)
			} else {
				ctlutil.UnauthenticatedError(ctx, w, r, "invalid api token", "request supplied invalid api token, denying")
			}
			return
		}

		// not supplying either is a valid use case, there are endpoints that allow anonymous access
		next.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(fn)
}

// --- accessors see ctxvalues ---
