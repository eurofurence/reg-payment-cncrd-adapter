package acceptance

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/concardis"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/config"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/paymentservice"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/web/app"
	"net/http/httptest"
)

// placing these here because they are package global

var (
	ts            *httptest.Server
	paymentMock   paymentservice.Mock
	concardisMock concardis.Mock
)

const tstConfigFile = "../resources/testconfig.yaml"

func tstSetup(configFilePath string) {
	tstSetupConfig(configFilePath)
	paymentMock = paymentservice.CreateMock()
	concardisMock = concardis.CreateMock()
	tstSetupHttpTestServer()
}

func tstSetupConfig(configFilePath string) {
	aulogging.SetupNoLoggerForTesting()
	config.LoadTestingConfigurationFromPathOrAbort(configFilePath)
}

func tstSetupHttpTestServer() {
	router := app.CreateRouter(context.Background())
	ts = httptest.NewServer(router)
}

func tstShutdown() {
	ts.Close()
	paymentMock.Reset()
	concardisMock.Reset()
}
