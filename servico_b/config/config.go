package config

import (
	"github.com/spf13/viper"
	"os"
)

type Conf struct {
	WeatherToken             string `mapstructure:"WEATHER_TOKEN"`
	OtelServiceName          string `mapstructure:"OTEL_SERVICE_NAME"`
	OtelExporterOtlpEndpoint string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
}

func LoadConfig(path string) (*Conf, error) {
	if os.Getenv("ENV") == "PROD" {
		return &Conf{
			WeatherToken:             os.Getenv("WEATHER_TOKEN"),
			OtelServiceName:          os.Getenv("OTEL_SERVICE_NAME"),
			OtelExporterOtlpEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		}, nil
	}
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	var cfg *Conf
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, nil
}
