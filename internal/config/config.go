package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-required:"true"`
	Database   `yaml:"database" env-required:"true"`
	HTTPServer `yaml:"http_server" env-required:"true"`
}

type Database struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     uint16 `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBname   string `yaml:"dbname" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

// 'must' prefix means panic
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH") // we can also pass '--config at the start'
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist in path: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config using cleanenv: %s", err)
	}

	return &cfg
}
