package config

import (
	"errors"
	"flag"
	"net/url"
	"os"
	"sort"
	"sync"

	aulogging "github.com/StephanHCB/go-autumn-logging"
	"gopkg.in/yaml.v2"
)

var (
	configurationData     *Application
	configurationLock     *sync.RWMutex
	configurationFilename string
	ecsLogging            bool
)

var (
	ErrorConfigArgumentMissing = errors.New("configuration file argument missing. Please specify using -config argument. Aborting")
	ErrorConfigFile            = errors.New("failed to read or parse configuration file. Aborting")
)

func init() {
	configurationData = &Application{Logging: LoggingConfig{Severity: "DEBUG"}}
	configurationLock = &sync.RWMutex{}

	flag.StringVar(&configurationFilename, "config", "", "config file path")
	flag.BoolVar(&ecsLogging, "ecs-json-logging", false, "switch to structured json logging")
}

// ParseCommandLineFlags is exposed separately so you can skip it for tests
func ParseCommandLineFlags() {
	flag.Parse()
}

func parseAndOverwriteConfig(yamlFile []byte, logPrintf func(format string, v ...interface{})) error {
	newConfigurationData := &Application{}
	err := yaml.UnmarshalStrict(yamlFile, newConfigurationData)
	if err != nil {
		logPrintf("failed to parse configuration file '%s': %v", configurationFilename, err)
		return err
	}

	setConfigurationDefaults(newConfigurationData)

	errs := url.Values{}
	validateServiceConfiguration(errs, newConfigurationData.Service)
	validateServerConfiguration(errs, newConfigurationData.Server)
	validateSecurityConfiguration(errs, newConfigurationData.Security)
	validateLoggingConfiguration(errs, newConfigurationData.Logging)
	validateInvoiceConfiguration(errs, newConfigurationData.Invoice)

	if len(errs) != 0 {
		var keys []string
		for key := range errs {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, k := range keys {
			key := k
			val := errs[k]
			logPrintf("configuration error: %s: %s", key, val[0])
		}
		return errors.New("configuration validation error")
	}

	configurationLock.Lock()
	defer configurationLock.Unlock()

	configurationData = newConfigurationData
	return nil
}

func loadConfiguration() error {
	yamlFile, err := os.ReadFile(configurationFilename)
	if err != nil {
		aulogging.Logger.NoCtx().Error().Printf("failed to load configuration file '%s': %v", configurationFilename, err)
		return err
	}
	err = parseAndOverwriteConfig(yamlFile, func(format string, v ...interface{}) {
		aulogging.Logger.NoCtx().Error().Printf(format, v...)
	})
	return err
}

// LoadTestingConfigurationFromPathOrAbort is for tests to set a hardcoded yaml configuration
func LoadTestingConfigurationFromPathOrAbort(configFilenameForTests string) {
	configurationFilename = configFilenameForTests
	if err := StartupLoadConfiguration(); err != nil {
		os.Exit(1)
	}
}

func StartupLoadConfiguration() error {
	aulogging.Logger.NoCtx().Info().Print("Reading configuration...")
	if configurationFilename == "" {
		aulogging.Logger.NoCtx().Error().Print("Configuration file argument missing. Please specify using -config argument. Aborting.")
		return ErrorConfigArgumentMissing
	}
	err := loadConfiguration()
	if err != nil {
		aulogging.Logger.NoCtx().Error().Print("Error reading or parsing configuration file. Aborting.")
		return ErrorConfigFile
	}
	return nil
}

func Configuration() *Application {
	configurationLock.RLock()
	defer configurationLock.RUnlock()
	return configurationData
}
