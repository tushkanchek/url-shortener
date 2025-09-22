package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	DBConfig `yaml:"db"`
	HTTPServer  `yaml:"http_server"`
}	


type DBConfig struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"dbname" env-default:"url-shortener"`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
}
type HTTPServer struct {
	Address     string      	`yaml:"address" env-default:"localhost:8083"`
	Timeout     time.Duration 	`yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration 	`yaml:"idle_timeout" env-default:"60s"`
	User 		string			`yaml:"user" env-required:"true"`
	Password 	string 			`yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`

}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}