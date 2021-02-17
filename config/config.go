package config

import (
	"github.com/caarlos0/env/v6"
)


type Config struct {
	PBIUrl   string `env:"PBI_URL"`
	Host   string `env:"HOST"`
	Port   int `env:"PORT"`
	Username   string `env:"USERNAME"`
	Password   string `env:"PASSWORD"`
	LogLevel   string `env:"LOG_LEVEL" envDefault:"info"`
	ConsumersCount   int `env:"CONSUMERS" envDefault:"1"`
}


func ReadConfig() (config *Config, err error) {
	config = new(Config)
	err = env.Parse(config)
	return config, err
}

