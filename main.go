package main

import (
	"github.com/eurofurence/reg-payment-cncrd-adapter/web"
)

func main() {
	web.StartWebserverAndNeverReturn()
}
