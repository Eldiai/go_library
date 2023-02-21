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
		//config.Port = 4000
		//config.Env = "development"
		//config.Db.Dsn = "postgres://postgres:postgres@localhost:5432/library?sslmode=disable"
		//config.Db.MaxOpenConns = 25
		//config.Db.MaxIdleConns = 25
		//config.Db.MaxIdleTime = "15m"
		//config.Smtp.Host = "smtp.office365.com"
		//config.Smtp.Port = 587
		//config.Smtp.Username = "" // change this
		//config.Smtp.Password = "" // change this
		//config.Smtp.Sender = ""   // change this
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
