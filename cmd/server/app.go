package app

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tomasz-trela/remitly-task/internal/database"
	"github.com/tomasz-trela/remitly-task/internal/handlers"
	"github.com/tomasz-trela/remitly-task/internal/seeders"
)

func Run() {
	log.Println("Starting server...")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	connStr := os.Getenv("DATABASE_URL")
	db, err := database.InitDB(connStr)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	for _, arg := range os.Args {
		if arg == "--seed" {
			log.Println("Seeding database...")
			seeders.SeedBanks()
			log.Println("Database seeded successfully.")
		}
	}

	r := handlers.NewRouter()

	http.ListenAndServe(":8080", r)
}
