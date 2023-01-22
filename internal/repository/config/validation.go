package config

import (
	"fmt"
	"net/url"
	"regexp"
)

func setConfigurationDefaults(c *Application) {
	if c.Server.Port == 0 {
		c.Server.Port = 8080
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

func validateServerConfiguration(errs url.Values, c ServerConfig) {
	checkIntValueRange(&errs, 1024, 65535, "server.port", int(c.Port))
	checkIntValueRange(&errs, 1, 300, "server.read_timeout_seconds", c.ReadTimeout)
	checkIntValueRange(&errs, 1, 300, "server.write_timeout_seconds", c.WriteTimeout)
	checkIntValueRange(&errs, 1, 300, "server.idle_timeout_seconds", c.IdleTimeout)
}

var allowedSeverities = []string{"DEBUG", "INFO", "WARN", "ERROR"}

func validateLoggingConfiguration(errs url.Values, c LoggingConfig) {
	if notInAllowedValues(allowedSeverities[:], c.Severity) {
		errs.Add("logging.severity", "must be one of DEBUG, INFO, WARN, ERROR")
	}
}

func validateSecurityConfiguration(errs url.Values, c SecurityConfig) {
	checkLength(&errs, 16, 256, "security.fixed.api", c.Fixed.Api)
	checkLength(&errs, 8, 64, "security.fixed.webhook", c.Fixed.Webhook)
}

const downstreamPattern = "^(|https?://.*[^/])$"

func validateServiceConfiguration(errs url.Values, c ServiceConfig) {
	if violatesPattern(downstreamPattern, c.PaymentService) {
		errs.Add("service.payment_service", "base url must be empty (enables in-memory simulator) or start with http:// or https:// and may not end in a /")
	}
	if violatesPattern(downstreamPattern, c.ConcardisDownstream) {
		errs.Add("service.concardis_downstream", "base url must be empty (enables local simulator) or start with http:// or https:// and may not end in a /")
	}
	if violatesPattern(downstreamPattern, c.PublicURL) {
		errs.Add("service.public_url", "public url must be empty or start with http:// or https:// and may not end in a /")
	}
	if c.ConcardisDownstream != "" && c.PublicURL != "" {
		errs.Add("service.public_url", "cannot set both public_url (for simulated paylinks) and concardis_downstream (to talk to actual api). Make up your mind!")
	}
	checkLength(&errs, 1, 256, "service.concardis_instance", c.ConcardisInstance)
	checkLength(&errs, 1, 256, "service.concardis_api_secret", c.ConcardisApiSecret)
}

func validateInvoiceConfiguration(errs url.Values, c InvoiceConfig) {
	checkLength(&errs, 1, 256, "invoice.title", c.Title)
	checkLength(&errs, 1, 256, "invoice.purpose", c.Purpose)
	checkLength(&errs, 1, 256, "invoice.description", c.Description)
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
