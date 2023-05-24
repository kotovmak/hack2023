package config

import (
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Tier       string `envconfig:"APP_TIER" default:"test"`
	Revision   string `envconfig:"REVISION" default:""`
	Version    string `envconfig:"VERSION" default:"0.0.1"`
	Host       string `envconfig:"APP_HOST" default:"localhost"`
	Port       int    `envconfig:"APP_PORT" default:"8081"`
	URL        string `envconfig:"DATABASE_URL" default:"root:password@tcp(localhost:3306)/hack2023?parseTime=true"`
	SigningKey string `envconfig:"SIGNING_KEY" default:"some_secret_key"`
}

var (
	cfg *Config
	mx  sync.RWMutex
)

func initConfig() (*Config, error) {
	godotenv.Load("/home/bitrix/www/hack2023/.env")
	cfg = &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Set main config
func Set(c Config) {
	mx.Lock()
	defer mx.Unlock()
	cfg = &c
}

// Get main config
func Get() Config {
	var err error
	mx.RLock()
	defer mx.RUnlock()
	if cfg == nil {
		cfg, err = initConfig()
		if err != nil {
			panic(err)
		}
	}
	return *cfg
}

// Reload main config
func Reload() Config {
	cfg, err := initConfig()
	if err != nil {
		return Get()
	}
	Set(*cfg)
	return Get()
}
