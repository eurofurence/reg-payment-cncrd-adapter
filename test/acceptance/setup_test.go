package acceptance

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/attendeeservice"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/concardis"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/config"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/mailservice"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/paymentservice"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/service/paymentlinksrv"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/web/app"
	"net/http/httptest"
	"time"
)

// placing these here because they are package global

var (
	ts            *httptest.Server
	attendeeMock  attendeeservice.Mock
	mailMock      mailservice.Mock
	paymentMock   paymentservice.Mock
	concardisMock concardis.Mock
)

const tstConfigFile = "../resources/testconfig.yaml"

const isoDateTimeFormat = "2006-01-02T15:04:05-07:00"

func tstMockNow() time.Time {
	mockTime, _ := time.Parse(isoDateTimeFormat, "2022-12-16T13:22:18+01:00")
	return mockTime
}

func tstSetup(configFilePath string) {
	tstSetupConfig(configFilePath)
	attendeeMock = attendeeservice.CreateMock()
	mailMock = mailservice.CreateMock()
	paymentMock = paymentservice.CreateMock()
	concardisMock = concardis.CreateMock()
	paymentlinksrv.NowFunc = tstMockNow
	tstSetupHttpTestServer()
}

func tstSetupConfig(configFilePath string) {
	aulogging.SetupNoLoggerForTesting()
	config.LoadTestingConfigurationFromPathOrAbort(configFilePath)
}

func tstSetupHttpTestServer() {
	router, _ := app.CreateRouter(context.Background())
	ts = httptest.NewServer(router)
}

func tstShutdown() {
	ts.Close()
	attendeeMock.Reset()
	mailMock.Reset()
	paymentMock.Reset()
	concardisMock.Reset()
}
