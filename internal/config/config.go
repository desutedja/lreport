package config

import (
	"os"

	"github.com/spf13/viper"
)

type config struct {
	LogLevel string
	Database *DatabaseConfig
	HTTPCfg  *HTTPConfig
}

var appConfig *config

func Load() error {
	viper.SetConfigName("application")
	env := os.Getenv("ENV")
	if env == "test" || env == "TEST" {
		viper.SetConfigName("test")
	}

	viper.SetConfigType("yaml")
	viper.AddConfigPath("./env/")
	viper.AddConfigPath("../env/")
	viper.AddConfigPath("../../env/")
	viper.AddConfigPath("../../../env/")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	viper.AutomaticEnv()

	appConfig = &config{
		LogLevel: getStringWithDefault("LOG_LEVEL", "warn"),
		Database: loadDatabaseConfig(),
		HTTPCfg:  loadHTTPConfig(),
	}

	return nil
}

func GetDatabaseConfig() *DatabaseConfig {
	return appConfig.Database
}

func GetHTTPConfig() *HTTPConfig {
	return appConfig.HTTPCfg
}
