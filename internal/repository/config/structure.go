package config

type (
	// Application is the root configuration type
	Application struct {
		Service  ServiceConfig  `yaml:"service"`
		Server   ServerConfig   `yaml:"server"`
		Logging  LoggingConfig  `yaml:"logging"`
		Security SecurityConfig `yaml:"security"`
		Invoice  InvoiceConfig  `yaml:"invoice"`
	}

	// ServiceConfig contains configuration values
	// for service related tasks. E.g. URLs to downstream services
	ServiceConfig struct {
		Name                string `yaml:"name"`
		PublicURL           string `yaml:"public_url"`           // my own public base url, without a trailing slash
		PaymentService      string `yaml:"payment_service"`      // base url, usually http://localhost:nnnn, will use in-memory-mock if unset
		ConcardisDownstream string `yaml:"concardis_downstream"` // base url, usually https://api.pay-link.eu, will use in-memory-mock if unset
		ConcardisInstance   string `yaml:"concardis_instance"`   // your instance name, required
		ConcardisApiSecret  string `yaml:"concardis_api_secret"` // your instance's api secret, required
	}

	// ServerConfig contains all values for http configuration
	ServerConfig struct {
		Address      string `yaml:"address"`
		Port         uint16 `yaml:"port"`
		ReadTimeout  int    `yaml:"read_timeout_seconds"`
		WriteTimeout int    `yaml:"write_timeout_seconds"`
		IdleTimeout  int    `yaml:"idle_timeout_seconds"`
	}

	// SecurityConfig configures everything related to incoming request security
	SecurityConfig struct {
		Fixed FixedTokenConfig `yaml:"fixed_token"`
		Cors  CorsConfig       `yaml:"cors"`
	}

	CorsConfig struct {
		DisableCors bool   `yaml:"disable"`
		AllowOrigin string `yaml:"allow_origin"`
	}

	FixedTokenConfig struct {
		Api     string `yaml:"api"`     // shared-secret for server-to-server backend authentication
		Webhook string `yaml:"webhook"` // shared-secret for the webhook coming in from concardis
	}

	// LoggingConfig configures logging
	LoggingConfig struct {
		Severity     string `yaml:"severity"`
		FullRequests bool   `yaml:"full_requests"`
	}

	// InvoiceConfig defines what the invoices should look like
	InvoiceConfig struct {
		Title       string `yaml:"title"`
		Description string `yaml:"description"`
		Purpose     string `yaml:"purpose"`
	}
)
