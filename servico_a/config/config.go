package config

import (
	"github.com/spf13/viper"
	"os"
)

// OTEL_SERVICE_NAME: service-a
// OTEL_EXPORTER_OTLP_ENDPOINT: otel:4317
type Conf struct {
	BaseUrl                  string `mapstructure:"BASE_URL"`
	OtelServiceName          string `mapstructure:"OTEL_SERVICE_NAME"`
	OtelExporterOtlpEndpoint string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
}

func LoadConfig(path string) (*Conf, error) {
	if os.Getenv("ENV") == "PROD" {
		return &Conf{
			BaseUrl:                  os.Getenv("BASE_URL"),
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
