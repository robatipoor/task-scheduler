package config

import (
	"time"

	"github.com/spf13/viper"
)

type Configure struct {
	Server struct {
		Schema string `mapstructure:"SCHEMA"`
		Host   string `mapstructure:"HOST"`
		Port   string `mapstructure:"PORT"`
	} `mapstructure:"SERVER"`

	Database struct {
		DSN string `mapstructure:"DSN"`
	} `mapstructure:"DATABASE"`

	Client struct {
		Timeout time.Duration `mapstructure:"TIMEOUT"`
	} `mapstructure:"CLIENT"`

	Scheduler struct {
		Duration time.Duration `mapstructure:"DURATION"`
	} `mapstructure:"SCHEDULER"`

	Master struct {
		Url string `mapstructure:"URL"`
	} `mapstructure:"MASTER"`
}

func LoadConfigure(path string) (*Configure, error) {
	viper.AddConfigPath(path)
	viper.SetEnvPrefix("WORKER")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	var config Configure
	if err := viper.ReadInConfig(); err != nil {
		return &config, err
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&config); err != nil {
		return &config, err
	}

	return &config, nil
}
