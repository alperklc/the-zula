package environment

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	LogLevel        string `env:"logLevel,required"`
	MongoURI        string `env:"mongoUri,required"`
	Environment     string `env:"environment,required"`
	Version         string `env:"version,required"`
	Port            int    `env:"port,required"`
	AuthClientId    string `env:"authClientId,required"`
	AuthKeyFilePath string `env:"authKeyFilePath,required"`
	AuthDomain      string `env:"authDomain,required"`
}

func Read() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to load .env file: %e", err)
	}

	cfg := Config{}
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse environment variables: %e", err)
	}

	fmt.Println("environment variables read.")

	return &cfg
}
