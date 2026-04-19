package main

import (
	"log"

	"service-core/internal/bootstrap"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, fallback to system env")
	}

	app := bootstrap.NewApp()
	app.Run()
}
