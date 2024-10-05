package config

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Database struct {
	Name     string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type HTTPServer struct {
	Address string `json:"address" yaml:"address"`
	Port    string `json:"port" yaml:"port"`
}

type Config struct {
	Env        string     `json:"env" yaml:"env" env-required:"true"`
	Database   Database   `json:"database" yaml:"database" env-required:"true"`
	HTTPServer HTTPServer `json:"http_server" yaml:"http_server"`
}

func MustGetConfig() Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		flags := flag.String("config", "", "path to config file")
		flag.Parse()

		configPath = *flags
		if configPath == "" {
			log.Fatal("config path not provided")
		}
	}

	_, err := os.Stat(configPath)
	if errors.Is(err, os.ErrNotExist) {
		log.Fatalf("config file %s not found", configPath)
	}

	var cfg Config

	fileReadError := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", fileReadError.Error())
	}
	return cfg
}
