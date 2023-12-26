package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env"`
	StoragePath string `yaml:"storage-path" env-required:"true"`
	Server      HTTPServer
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	TimeOut     time.Duration `yaml:"timeout"`
	IdleTimeOut time.Duration `yaml:"idle-timeout"`
}

func MustLoad() *Config {
	configPath := "C:/Users/maus1/Desktop/url-shortener/config/local.yml" // configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("Config path is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	return &cfg
}
