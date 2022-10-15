package config

import (
	"github.com/eurofurence/reg-payment-cncrd-adapter/docs"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseAndOverwriteConfigInvalidYamlSyntax(t *testing.T) {
	docs.Description("check that a yaml with a syntax error leads to a parse error")
	invalidYaml := `# invalid yaml
security:
    disable_cors: true # indented wrong
  fixed_token:
    api: # no value
`
	err := parseAndOverwriteConfig([]byte(invalidYaml))
	require.NotNil(t, err, "expected an error")
}

func TestParseAndOverwriteConfigUnexpectedFields(t *testing.T) {
	docs.Description("check that a yaml with unexpected fields leads to a parse error")
	invalidYaml := `# yaml with model mismatches
serval:
  port: 8088
cheetah:
  speed: '60 mph'
`
	err := parseAndOverwriteConfig([]byte(invalidYaml))
	require.NotNil(t, err, "expected an error")
}

func TestParseAndOverwriteConfigValidationErrors1(t *testing.T) {
	docs.Description("check that a yaml with validation errors leads to an error")
	wrongConfigYaml := `# yaml with validation errors
server:
  port: 14
logging:
  severity: FELINE
`
	err := parseAndOverwriteConfig([]byte(wrongConfigYaml))
	require.NotNil(t, err, "expected an error")
	require.Equal(t, err.Error(), "configuration validation error", "unexpected error message")
}

func TestParseAndOverwriteDefaults(t *testing.T) {
	docs.Description("check that a minimal yaml leads to all defaults being set")
	minimalYaml := `# yaml with minimal settings
security:
  fixed_token:
    api: 'fixed-testing-token-abc'
    webhook: 'fixed-webhook-token-abc'
`
	err := parseAndOverwriteConfig([]byte(minimalYaml))
	require.Nil(t, err, "expected no error")
	require.Equal(t, uint16(8080), Configuration().Server.Port, "unexpected value for server.port")
	require.Equal(t, "INFO", Configuration().Logging.Severity, "unexpected value for logging.severity")
}
