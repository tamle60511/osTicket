package main

import (
	"ecommerce/config"
	"ecommerce/internal/api"
	"log"
)

func main() {
	cfg, err := config.SetupEnv()
	if err != nil {
		log.Printf("Server Env not config")
	}
	api.StartServer(cfg)
}
