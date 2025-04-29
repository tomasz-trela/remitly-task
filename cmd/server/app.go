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
	"github.com/tomasz-trela/remitly-task/internal/models"
	"github.com/tomasz-trela/remitly-task/internal/parsers"
	"github.com/tomasz-trela/remitly-task/internal/repository"
	"github.com/tomasz-trela/remitly-task/internal/seeders"
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

	for _, arg := range os.Args {
		if arg == "--seed" {
			log.Println("Seeding database...")
			seeders.SeedBanks()
			log.Println("Database seeded successfully.")
		}
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/swift-codes/{swiftCode}", getBankDataBySwift)
	r.Get("/swift-codes/country/{countryISO2code}", getBanksDataByCountryISO2)
	r.Post("/swift-codes", createBank)
	r.Delete("/swift-codes/{swiftCode}", DeleteSwiftCode)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	parsers.LoadSwiftRecords()
	http.ListenAndServe(":8080", r)
}

func getBankDataBySwift(w http.ResponseWriter, r *http.Request) {
	swiftCode := chi.URLParam(r, "swiftCode")

	swiftCodeResponse, err := repository.GetBankCodeAndBranchesBySwift(swiftCode)
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

func getBanksDataByCountryISO2(w http.ResponseWriter, r *http.Request) {
	countryISO2 := chi.URLParam(r, "countryISO2code")

	banks, err := repository.GetBanksByISO2(countryISO2)
	if err != nil {
		http.Error(w, "Error retrieving bank data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(banks)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
	}
}

func createBank(w http.ResponseWriter, r *http.Request) {
	var swiftCode models.SwiftCode
	err := json.NewDecoder(r.Body).Decode(&swiftCode)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	err = repository.InsertSwiftCode(&swiftCode)
	if err != nil {
		http.Error(w, "Error inserting bank data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Bank created successfully"))
}

func DeleteSwiftCode(w http.ResponseWriter, r *http.Request) {
	swiftCode := chi.URLParam(r, "swiftCode")

	err := repository.DeleteSwiftCode(swiftCode)
	if err != nil {
		http.Error(w, "Error deleting bank data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
