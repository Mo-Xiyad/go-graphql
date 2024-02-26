package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type database struct {
	URL string
}

type jwt struct {
	Secret string
	Issuer string
}

type Config struct {
	Database database
	JWT      jwt
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func New() *Config {
	return &Config{
		Database: database{
			URL: fmt.Sprint(os.Getenv("DSN")),
		},
		JWT: jwt{
			Secret: fmt.Sprint(os.Getenv("JWT_SECRET")),
			Issuer: fmt.Sprint(os.Getenv("ISSUER")),
		},
	}
}
