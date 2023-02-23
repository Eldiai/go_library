package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Db struct {
		Dsn          string `json:"dsn" yaml:"dsn"`
		MaxOpenConns int    `json:"maxOpenConns" yaml:"maxOpenConns"`
		MaxIdleConns int    `json:"maxIdleConns" yaml:"maxIdleConns"`
		MaxIdleTime  string `json:"maxIdleTime" yaml:"maxIdleTime"`
	}

	Smtp struct {
		Host     string `json:"host" yaml:"host"`
		Port     int    `json:"port" yaml:"port"`
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
		Sender   string `json:"sender" yaml:"sender"`
	}

	Config struct {
		Port string `json:"port" yaml:"port"`
		Env  string `json:"env" yaml:"env"`
		Db   *Db    `json:"db" yaml:"db"`
		Smtp *Smtp  `json:"smtp" yaml:"smtp"`
	}
)

var (
	once   sync.Once
	config = new(Config)
)

func GetConfig() *Config {
	once.Do(func() {
		if err := cleanenv.ReadConfig("config/config.yml", config); err != nil {
			panic(err)
		}
	})

	return config
}
