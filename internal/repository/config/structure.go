package config

type serverConfig struct {
	Address      string `yaml:"address"`
	Port         uint16 `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout_seconds"`
	WriteTimeout int    `yaml:"write_timeout_seconds"`
	IdleTimeout  int    `yaml:"idle_timeout_seconds"`
}

type downstreamConfig struct {
	PaymentService      string `yaml:"payment_service"`      // base url, usually http://localhost:nnnn, will use in-memory-mock if unset
	ConcardisDownstream string `yaml:"concardis_downstream"` // base url, usually https://api.pay-link.eu, will use in-memory-mock if unset
	ConcardisInstance   string `yaml:"concardis_instance"`   // your instance name, required
	ConcardisApiSecret  string `yaml:"concardis_api_secret"` // your instance's api secret, required
}

type loggingConfig struct {
	Severity string `yaml:"severity"`
}

type fixedTokenConfig struct {
	Api     string `yaml:"api"`     // shared-secret for server-to-server backend authentication
	Webhook string `yaml:"webhook"` // shared-secret for the webhook coming in from concardis
}

type securityConfig struct {
	Fixed       fixedTokenConfig `yaml:"fixed_token"`
	DisableCors bool             `yaml:"disable_cors"`
}

type conf struct {
	Server     serverConfig     `yaml:"server"`
	Logging    loggingConfig    `yaml:"logging"`
	Security   securityConfig   `yaml:"security"`
	Downstream downstreamConfig `yaml:"downstream"`
}
