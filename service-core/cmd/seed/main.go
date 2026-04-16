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

	conn := database.NewConnection(cfg)
	defer conn.Close()

	database.RunSeed(conn)

	log.Println("seed completed")
}
