package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Appconfig struct {
	ServerPort string
	DSN        string
	AppSecret  string
}

func SetupEnv() (c Appconfig, err error) {
	if os.Getenv("APP_ENV") == "dev" {
		godotenv.Load()
	}
	httpPort := os.Getenv("HTTP_PORT")
	if len(httpPort) < 1 {
		return Appconfig{}, errors.New("HTTP PORT not config")
	}
	Dsn := os.Getenv("DSN")
	if len(Dsn) < 1 {
		return Appconfig{}, errors.New("DSN not config")
	}
	appSecret := os.Getenv("APP_SECRET")
	if len(appSecret) < 1 {
		return Appconfig{}, errors.New("APP_SECRET not config")
	}
	return Appconfig{ServerPort: httpPort, DSN: Dsn, AppSecret: appSecret}, nil
}
