package config

import (
	"errors"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)


type Config struct {
	DB DB `yaml:"db"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Mode string `yaml:"mode"`
}

type HTTPServer struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port int `yaml:"port" env-default:"8000"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"120"`
	ReadTimeout time.Duration `yaml:"read_timeout" env-default:"5"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"5"`
}

type DB struct {
	User string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Name string `yaml:"name" env-required:"true"`
	Host string `yaml:"host" env-default:"localhost"`
	Port int `yaml:"port" env-required:"true"`
}

func Load(configPath string) (*Config, error) {
	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		return nil, errors.New("failed to load configuration file: " + err.Error())
	}
	return &config, nil
}