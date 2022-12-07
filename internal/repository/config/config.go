// configuration management using a yaml configuration file
// You must have called LoadConfiguration() or otherwise set up the configuration before you can use these.
package config

import (
	"fmt"
	"time"
)

func UseEcsLogging() bool {
	return ecsLogging
}

func ServerAddr() string {
	c := Configuration()
	return fmt.Sprintf("%s:%d", c.Server.Address, c.Server.Port)
}

func ServerReadTimeout() time.Duration {
	return time.Second * time.Duration(Configuration().Server.ReadTimeout)
}

func ServerWriteTimeout() time.Duration {
	return time.Second * time.Duration(Configuration().Server.WriteTimeout)
}

func ServicePublicURL() string {
	return Configuration().Service.PublicURL
}

func ServerIdleTimeout() time.Duration {
	return time.Second * time.Duration(Configuration().Server.IdleTimeout)
}

func LoggingSeverity() string {
	return Configuration().Logging.Severity
}

func FixedApiToken() string {
	return Configuration().Security.Fixed.Api
}

func IsCorsDisabled() bool {
	return Configuration().Security.Cors.DisableCors
}

func PaymentServiceBaseUrl() string {
	return Configuration().Service.PaymentService
}

func ConcardisDownstreamBaseUrl() string {
	return Configuration().Service.ConcardisDownstream
}

func ConcardisInstanceName() string {
	return Configuration().Service.ConcardisInstance
}

func ConcardisInstanceApiSecret() string {
	return Configuration().Service.ConcardisApiSecret
}

func WebhookSecret() string {
	return Configuration().Security.Fixed.Webhook
}
