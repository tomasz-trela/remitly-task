package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/tomasz-trela/remitly-task/internal/database"
	"github.com/tomasz-trela/remitly-task/internal/parsers"
	"github.com/tomasz-trela/remitly-task/internal/repository"
)

func Run() {
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

	// seeders.SeedBanks()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/swift-codes/{swiftCode}", getBankDataBySwift)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	parsers.LoadSwiftRecords()
	http.ListenAndServe(":8080", r)
}

func getBankDataBySwift(w http.ResponseWriter, r *http.Request) {
	swiftCode := chi.URLParam(r, "swiftCode")

	swiftCodeResponse, err := repository.GetBankCodeBySwift(swiftCode)
	if err != nil {
		http.Error(w, "Error retrieving bank data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(swiftCodeResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
	}
}
