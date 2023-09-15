package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type (
	Config struct {
		LogLevel   string `env:"LOG_LEVEL" env-default:"debug"`
		HTTP       HTTP
		PostgreSQL DBConfig
	}
	HTTP struct {
		IsHTTPS        bool   `env:"IS_HTTPS" env-required:"true"`
		ServerCertPath string `env:"SERVER_CERT_PATH" env-required:"true"`
		PrivateKeyPath string `env:"PRIVATE_KEY_PATH" env-required:"true"`
	}
	DBConfig struct {
		Port     string `env:"PSQL_PORT" env-required:"true"`
		Host     string `env:"PSQL_HOST" env-required:"true"`
		Username string `env:"PSQL_USERNAME" env-required:"true"`
		Database string `env:"PSQL_DATABASE" env-required:"true"`
		Password string `env:"PSQL_PASSWORD" env-required:"true"`
	}
)

// nolint
var instance *Config

// nolint
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadEnv(instance); err != nil {
			//if err := cleanenv.ReadConfig(".env", instance); err != nil {
			helpText := "error read env"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Println(help)
			log.Fatal(err)
		}
	})

	return instance
}
