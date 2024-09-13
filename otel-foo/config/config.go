package config

import (
	"github.com/spf13/viper"
)

func InitConf() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			panic("Config.yaml file not found.")
		default:
			panic("Failed to load config.yaml.")
		}
	}
	return nil
}

type httpConfig struct {
	Port string `mapstructure:"port"`
}

type barServerConfig struct {
	Host string `mapstructure:"host"`
}

type otelCollectorConfig struct {
	Host string `mapstructure:"host"`
}

func GetHttpConfig() httpConfig {
	return httpConfig{
		Port: viper.GetString("http.port"),
	}
}

func GetBarServerConfig() barServerConfig {
	return barServerConfig{
		Host: viper.GetString("bar.host"),
	}
}

func GetOtelCollectorConfig() otelCollectorConfig {
	return otelCollectorConfig{
		Host: viper.GetString("otel_collector.host"),
	}
}
