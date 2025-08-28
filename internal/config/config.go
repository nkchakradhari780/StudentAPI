package config

import (
	"flag"
	"log"
	"os"
	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Addr string `yaml:"address" env-required:"true"`
}
type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HttpServer  HttpServer `yaml:"http_server" `
}

func MustLoad() *Config {
	var configPath string

	configPath =  os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "config/local.yaml", "Path to the config file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path must be specified")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err.Error())
	}

	return &cfg

}