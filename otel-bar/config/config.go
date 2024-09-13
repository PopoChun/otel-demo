package config

import "github.com/spf13/viper"

type Config struct {
	HttpPort          string `mapstructure:"HTTP_PORT"`
	OtelCollectorHost string `mapstructure:"OTEL_COLLECTOR_HOST"`
}

func InitConf() (config Config, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		panic(err.(viper.ConfigFileNotFoundError))
	}

	err = viper.Unmarshal(&config)
	return
}
