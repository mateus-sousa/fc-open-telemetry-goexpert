package config

import (
	"github.com/spf13/viper"
	"os"
)

type Conf struct {
	WeatherToken string `mapstructure:"WEATHER_TOKEN"`
	ExporterUrl  string `mapstructure:"EXPORTER_URL"`
}

func LoadConfig(path string) (*Conf, error) {
	if os.Getenv("ENV") == "PROD" {
		return &Conf{
			WeatherToken: os.Getenv("WEATHER_TOKEN"),
			ExporterUrl:  os.Getenv("EXPORTER_URL"),
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
