package config

import (
	"fmt"
	"testing"

	"github.com/eurofurence/reg-payment-cncrd-adapter/docs"
	"github.com/stretchr/testify/require"
)

var recording []string

func tstLogRecorder(format string, v ...interface{}) {
	recording = append(recording, fmt.Sprintf(format, v...))
}

func TestParseAndOverwriteConfigInvalidYamlSyntax(t *testing.T) {
	docs.Description("check that a yaml with a syntax error leads to a parse error")
	invalidYaml := `# invalid yaml
security:
    disable_cors: true # indented wrong
  fixed_token:
    api: # no value
`
	recording = make([]string, 0)
	err := parseAndOverwriteConfig([]byte(invalidYaml), tstLogRecorder)
	require.NotNil(t, err, "expected an error")
	require.Equal(t, 1, len(recording))
	require.Contains(t, recording[0], "failed to parse configuration file")
}

func TestParseAndOverwriteConfigUnexpectedFields(t *testing.T) {
	docs.Description("check that a yaml with unexpected fields leads to a parse error")
	invalidYaml := `# yaml with model mismatches
serval:
  port: 8088
cheetah:
  speed: '60 mph'
`
	recording = make([]string, 0)
	err := parseAndOverwriteConfig([]byte(invalidYaml), tstLogRecorder)
	require.NotNil(t, err, "expected an error")
	require.Equal(t, 1, len(recording))
	require.Contains(t, recording[0], "failed to parse configuration file")
}

func TestParseAndOverwriteConfigValidationErrors1(t *testing.T) {
	docs.Description("check that a yaml with validation errors leads to an error")
	wrongConfigYaml := `# yaml with validation errors
service:
  public_url: 'https://invalid.has.trailing.slash/'
  payment_service: 'also not a valid url'
  concardis_downstream: 'another invalid url'
server:
  port: 14
logging:
  severity: FELINE
`
	recording = make([]string, 0)
	err := parseAndOverwriteConfig([]byte(wrongConfigYaml), tstLogRecorder)
	require.NotNil(t, err, "expected an error")
	require.Equal(t, err.Error(), "configuration validation error", "unexpected error message")
	require.EqualValues(t, []string{
		"configuration error: database.use: must be one of mysql, inmemory",
		"configuration error: invoice.description: invoice.description field must be at least 1 and at most 256 characters long",
		"configuration error: invoice.purpose: invoice.purpose field must be at least 1 and at most 256 characters long",
		"configuration error: invoice.title: invoice.title field must be at least 1 and at most 256 characters long",
		"configuration error: logging.severity: must be one of DEBUG, INFO, WARN, ERROR",
		"configuration error: security.fixed.api: security.fixed.api field must be at least 16 and at most 256 characters long",
		"configuration error: security.fixed.webhook: security.fixed.webhook field must be at least 8 and at most 64 characters long",
		"configuration error: server.port: server.port field must be an integer at least 1024 and at most 65535",
		"configuration error: service.concardis_api_secret: service.concardis_api_secret field must be at least 1 and at most 256 characters long",
		"configuration error: service.concardis_downstream: base url must be empty (enables local simulator) or start with http:// or https:// and may not end in a /",
		"configuration error: service.concardis_instance: service.concardis_instance field must be at least 1 and at most 256 characters long",
		"configuration error: service.payment_service: base url must be empty (enables in-memory simulator) or start with http:// or https:// and may not end in a /",
		"configuration error: service.public_url: public url must be empty or start with http:// or https:// and may not end in a /",
	}, recording)
}

func TestParseAndOverwriteDefaults(t *testing.T) {
	docs.Description("check that a minimal yaml leads to all defaults being set")
	minimalYaml := `# yaml with minimal settings
security:
  fixed_token:
    api: 'fixed-testing-token-abc'
    webhook: 'fixed-webhook-token-abc'
database:
  use: inmemory
service:
  concardis_instance: 'my-demo-instance'
  concardis_api_secret: 'my-demo-secret'
invoice:
  title: 'demo title'
  description: 'demo description'
  purpose: 'demo purpose'
`
	recording = make([]string, 0)
	err := parseAndOverwriteConfig([]byte(minimalYaml), tstLogRecorder)
	require.Nil(t, err, "expected no error")
	require.Equal(t, uint16(8080), Configuration().Server.Port, "unexpected value for server.port")
	require.Equal(t, "INFO", Configuration().Logging.Severity, "unexpected value for logging.severity")
}
