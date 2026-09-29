package config

import (
	"flag"
	"fmt"
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

	configPath := resolveConfigPath()

	fmt.Println("config path: ", configPath)

	if configPath == "" {
		slog.Error("config path is empty")
		os.Exit(1)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("config file does not exist at path", "configPath", configPath)
		os.Exit(1)
	}

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		slog.Error("cannot read config", "err", err.Error())
		os.Exit(1)
	}

	return &config
}

func resolveConfigPath() string {
	var configPath string

	flags := flag.String("config", "", "Pathto the config file")
	flag.Parse()
	configPath = *flags

	if configPath != "" {
		slog.Info("config is form flags")
		return configPath
	}

	configPath = os.Getenv("CONFIG_PATH")

	if configPath != "" {
		slog.Info("config is form env")
		return configPath
	}

	configPath = "./config/local.config.yaml"

	slog.Info("config is form local yaml")

	return configPath
}
