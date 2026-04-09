package main

import (
	"log"
	database "service-core/internal/infra/db"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, fallback to system env")
	}

	cfg := database.LoadConfig()

	database.RunMigrations(cfg)

	log.Println("migration completed")
}
