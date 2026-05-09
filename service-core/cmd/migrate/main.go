package main

import (
	"log"

	"service-core/internal/bootstrap"
	database "service-core/internal/infra/db"
	"service-core/internal/infra/storage"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, fallback to system env")
	}

	cfg := bootstrap.LoadConfig()

	infra, err := bootstrap.NewInfra(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer infra.Close()

	if err := database.RunMigration(cfg.DB); err != nil {
		log.Fatal(err)
	}
	log.Println("database migration complete")

	if err := storage.RunMigration(infra.StorageProvider); err != nil {
		log.Fatal(err)
	}
	log.Println("storage migration complete")

	log.Println("migration complete")
}
