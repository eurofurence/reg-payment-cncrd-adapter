package infoctl

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Create(server chi.Router) {
	server.Get("/", healthHandler)
	server.Get("/info/health", healthHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "OK")
}
