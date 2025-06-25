package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"development"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Addres      string        `yaml:"address" env-default:"localhost8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	//Создаем конфигурационный путь
	congifPath := os.Getenv("CONFIG_PATH")
	if congifPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	//Проверяем существование конфигурационного файла
	if _, err := os.Stat(congifPath); os.IsNotExist(err) {
		log.Fatalf("logger file does not exist: %s", congifPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(congifPath, &cfg); err != nil {
		log.Fatalf("error reading logger: %s", err)
	}

	return &cfg

}
