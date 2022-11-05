package main

import (
	"balance-microservice/internal/app/config"
	"balance-microservice/internal/app/server"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

func main() {
	cfg := config.NewConfig()
	b, readError := os.ReadFile("./configs/application.yaml")
	if readError != nil {
		log.Fatal("Error on reading config")
	}

	err := yaml.Unmarshal(b, &cfg)
	if err != nil {
		return
	}

	if err := server.Start(cfg); err != nil {
		log.Fatal("Error on starting server")
	}

}
