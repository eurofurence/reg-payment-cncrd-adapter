package concardis

import (
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/config"
)

var activeInstance ConcardisDownstream

func Create() (err error) {
	if config.ConcardisDownstreamBaseUrl() != "" {
		activeInstance, err = newClient()
		return err
	} else {
		aulogging.Logger.NoCtx().Warn().Print("service.concardis_downstream not configured. Using in-memory simulator for concardis downstream (not useful for production!)")
		activeInstance = newMock()
		return nil
	}
}

func CreateMock() Mock {
	instance := newMock()
	activeInstance = instance
	return instance
}

func Get() ConcardisDownstream {
	return activeInstance
}
