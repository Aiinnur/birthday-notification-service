package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	PostgresURL string     `yaml:"postgres_url" env-required:"true"`
	Server      HTTPServer `yaml:"http_server"`
	SmtpServer  SMTPServer `yaml:"smtp_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type SMTPServer struct {
	Address  string `yaml:"smtpAddress" env-required:"true"`
	Port     string `yaml:"smtpPort" env-default:"587"`
	User     string `yaml:"smtpUser" env-required:"true"`
	Password string `yaml:"smtpPassword" env-required:"true"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		log.Fatal("Config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal("Config file does not exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatal("Failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
