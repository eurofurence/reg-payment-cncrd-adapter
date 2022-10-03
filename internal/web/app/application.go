package app

import (
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/config"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/paymentservice"
)

type Application interface {
	Run() int
}

type Impl struct{}

func New() Application {
	return &Impl{}
}

func (i *Impl) Run() int {
	config.ParseCommandLineFlags()
	setupLogging("attendee-service", config.UseEcsLogging())

	if err := config.StartupLoadConfiguration(); err != nil {
		return 1
	}
	setLoglevel(config.LoggingSeverity())

	if err := paymentservice.Create(); err != nil {
		return 1
	}

	if err := runServerWithGracefulShutdown(); err != nil {
		return 2
	}

	return 0
}
