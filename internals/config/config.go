package config

import (
	"errors"
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "config path")
}

type Config struct {
	EnvLevel    int    `yaml:"env_level" env-default:"1"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	GRPC        GRPC   `yaml:"grpc"`
}

type GRPC struct {
	Port    int           `yaml:"port" env-default:"44044"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

func MustLoadByString(path string) *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}

func MustLoad() *Config {
	flag.Parse()

	if err := fetchConfigPath(); err != nil {
		panic(err)
	}

	return MustLoadByString(configPath)
}

func fetchConfigPath() error {
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
		if configPath == "" {
			return errors.New("config path is required")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return errors.New("config path does not exist")
	}

	return nil
}
