package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tomasz-trela/remitly-task/internal/models"
	"github.com/tomasz-trela/remitly-task/internal/repository"
)

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.MessageResponse{Message: message})
}

func writeJSONMessage(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.MessageResponse{Message: message})
}

func GetBankDataBySwift(w http.ResponseWriter, r *http.Request) {
	swiftCode := chi.URLParam(r, "swiftCode")

	swiftCodeResponse, err := repository.GetBankCodeAndBranchesBySwift(swiftCode)
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "No bank found with the given swift code")
		return
	}
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Error retrieving bank data")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(swiftCodeResponse)
}

func GetBanksDataByCountryISO2(w http.ResponseWriter, r *http.Request) {
	countryISO2 := chi.URLParam(r, "countryISO2code")

	banks, err := repository.GetBanksByISO2(countryISO2)
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "No banks found for the given country ISO2 code")
		return
	}
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Error retrieving bank data")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banks)
}

func CreateBank(w http.ResponseWriter, r *http.Request) {
	var swiftCode models.SwiftCode
	err := json.NewDecoder(r.Body).Decode(&swiftCode)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	rowsAffected, err := repository.InsertSwiftCode(&swiftCode)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Error inserting bank data")
		return
	}

	if rowsAffected == 0 {
		writeJSONError(w, http.StatusConflict, fmt.Sprintf("Swift code %s already exists", swiftCode.SwiftCode))
		return
	}

	writeJSONMessage(w, http.StatusCreated, "Bank created successfully")
}

func DeleteSwiftCode(w http.ResponseWriter, r *http.Request) {
	swiftCode := chi.URLParam(r, "swiftCode")

	rowsAffected, err := repository.DeleteSwiftCode(swiftCode)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Error deleting bank data")
		return
	}
	if rowsAffected == 0 {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("No bank found with swift code %s", swiftCode))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
