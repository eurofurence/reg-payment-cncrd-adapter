package main

import (
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/web/app"
	"os"
)

func main() {
	os.Exit(app.New().Run())
}
