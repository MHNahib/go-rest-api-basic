package config

import (
	"flag"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Address string `yaml:"address"`
}

type Database struct {
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}

type Config struct {
	Env string `yaml:"env"`

	Server   HttpServer `yaml:"server"`
	Database Database   `yaml:"database"`
}

func MountConfig() *Config {
	var config Config

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "Pathto the config file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			slog.Error("config path is empty")
			os.Exit(1)
		}

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			slog.Error("config file does not exist at path: ", "configPath", configPath)
			os.Exit(1)
		}

		if err := cleanenv.ReadConfig(configPath, &config); err != nil {
			slog.Error("cannot read config", "err", err.Error())
			os.Exit(1)
		}

	}

	return &config
}
