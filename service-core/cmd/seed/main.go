package main

import (
	"log"

	"service-core/internal/bootstrap"
	database "service-core/internal/infra/db"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("\x1b[0;94;49mseed start\x1b[0;39;49m")

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, fallback to system env")
	}

	cfg := bootstrap.LoadConfig()
	infra, err := bootstrap.NewInfra(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer infra.Close()

	if err := database.RunSeed(infra.DB); err != nil {
		log.Fatal(err)
	}

	log.Println("\x1b[0;94;49mseed complete\x1b[0;39;49m")
}
