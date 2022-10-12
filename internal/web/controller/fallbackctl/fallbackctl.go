package fallbackctl

import (
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/web/util/ctlutil"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Create(server chi.Router) {
	server.HandleFunc("/*", fallbackErrorHandler)
}

func fallbackErrorHandler(w http.ResponseWriter, r *http.Request) {
	ctlutil.ErrorHandler(r.Context(), w, r, "not.found", http.StatusNotFound, nil)
}
