package config

import (
	"fmt"
	"net/url"
	"regexp"
)

func setConfigurationDefaults(c *conf) {
	if c.Server.Port == "" {
		c.Server.Port = "8080"
	}
	if c.Server.ReadTimeout <= 0 {
		c.Server.ReadTimeout = 5
	}
	if c.Server.WriteTimeout <= 0 {
		c.Server.WriteTimeout = 5
	}
	if c.Server.IdleTimeout <= 0 {
		c.Server.IdleTimeout = 5
	}
	if c.Logging.Severity == "" {
		c.Logging.Severity = "INFO"
	}
}

const portPattern = "^[1-9][0-9]{0,4}$"

func validateServerConfiguration(errs url.Values, c serverConfig) {
	if violatesPattern(portPattern, c.Port) {
		errs.Add("server.port", "must be a number between 1 and 65535")
	}
	checkIntValueRange(&errs, 1, 300, "server.read_timeout_seconds", c.ReadTimeout)
	checkIntValueRange(&errs, 1, 300, "server.write_timeout_seconds", c.WriteTimeout)
	checkIntValueRange(&errs, 1, 300, "server.idle_timeout_seconds", c.IdleTimeout)
}

var allowedSeverities = []string{"DEBUG", "INFO", "WARN", "ERROR"}

func validateLoggingConfiguration(errs url.Values, c loggingConfig) {
	if notInAllowedValues(allowedSeverities[:], c.Severity) {
		errs.Add("logging.severity", "must be one of DEBUG, INFO, WARN, ERROR")
	}
}

func validateSecurityConfiguration(errs url.Values, c securityConfig) {
	checkLength(&errs, 16, 256, "security.fixed.api", c.Fixed.Api)
}

const downstreamPattern = "^(|https?://.*[^/])$"

func validateDownstreamConfiguration(errs url.Values, c downstreamConfig) {
	if violatesPattern(downstreamPattern, c.PaymentService) {
		errs.Add("downstream.payment_service", "base url must be empty (enables in-memory simulator) or start with http:// or https:// and may not end in a /")
	}
}

// -- helpers

func violatesPattern(pattern string, value string) bool {
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return true
	}
	return !matched
}

func checkLength(errs *url.Values, min int, max int, key string, value string) {
	if len(value) < min || len(value) > max {
		errs.Add(key, fmt.Sprintf("%s field must be at least %d and at most %d characters long", key, min, max))
	}
}

func checkIntValueRange(errs *url.Values, min int, max int, key string, value int) {
	if value < min || value > max {
		errs.Add(key, fmt.Sprintf("%s field must be an integer at least %d and at most %d", key, min, max))
	}
}

func notInAllowedValues(allowed []string, value string) bool {
	for _, v := range allowed {
		if v == value {
			return false
		}
	}
	return true
}
