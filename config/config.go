package config

import (
	"encoding/json"
	"os"
	"sync"
)

type (
	Config struct {
		Port string `json:"port"`
		Env  string `json:"env"`
		Db   struct {
			Dsn          string `json:"dsn"`
			MaxOpenConns int    `json:"maxOpenConns"`
			MaxIdleConns int    `json:"maxIdleConns"`
			MaxIdleTime  string `json:"maxIdleTime"`
		} `json:"db"`
		Smtp struct {
			Host     string `json:"host"`
			Port     int    `json:"port"`
			Username string `json:"username"`
			Password string `json:"password"`
			Sender   string `json:"sender"`
		} `json:"smtp"`
	}
)

var (
	once   sync.Once
	config = new(Config)
)

func GetConfig() *Config {
	once.Do(func() {
		b, err := os.ReadFile("config/default.json")
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(b, config)
		if err != nil {
			panic(err)
		}
	})

	return config
}
