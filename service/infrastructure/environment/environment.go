package environment

import (
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	LogLevel                 string `env:"LOG_LEVEL"`
	MongoURI                 string `env:"MONGO_URI"`
	Environment              string `env:"ENVIRONMENT"`
	Version                  string `env:"VERSION"`
	Port                     int    `env:"PORT,required"`
	AuthClientId             string `env:"AUTH_CLIENT_ID"`
	AuthKeyFilePath          string `env:"AUTH_KEY_FILE_PATH"`
	AuthDomain               string `env:"AUTH_DOMAIN"`
	AuthServiceAccountUser   string `env:"AUTH_SERVICE_ACCOUNT_USER"`
	AuthServiceAccountSecret string `env:"AUTH_SERVICE_ACCOUNT_SECRET"`
	RabbitMqUser             string `env:"RABBIT_MQ_USER"`
	RabbitMqPassword         string `env:"RABBIT_MQ_PASSWORD"`
	RabbitMqUri              string `env:"RABBIT_MQ_URI"`
}

func Read() *Config {
	// Try to load the .env file (optional)
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("error loading .env file: %v", err)
		} else {
			fmt.Println(".env file loaded successfully.")
		}
	} else {
		fmt.Println("No .env file found, reading environment variables from system.")
	}

	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse environment variables: %e", err)
	}

	fmt.Println("environment variables read.")

	return &cfg
}
