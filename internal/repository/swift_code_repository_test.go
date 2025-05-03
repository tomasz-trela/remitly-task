package repository

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tomasz-trela/remitly-task/internal/database"
	"github.com/tomasz-trela/remitly-task/internal/models"
)

func TestInsertSwiftCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	database.DB = db

	swiftCode := &models.SwiftCode{
		SwiftCode:     "BPHKPLPK",
		IsHeadquarter: false,
		BankName:      "Bank BPH",
		CountryISO2:   "PL",
		CountryName:   "Poland",
		Address:       "ul. Towarowa 25A",
	}

	mock.ExpectExec("^INSERT INTO countries").
		WithArgs(swiftCode.CountryISO2, swiftCode.CountryName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("^INSERT INTO banks").
		WithArgs(swiftCode.SwiftCode, swiftCode.IsHeadquarter, swiftCode.BankName, swiftCode.Address, swiftCode.CountryISO2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = InsertSwiftCode(swiftCode)
	if err != nil {
		t.Errorf("Error not expected: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestInsertSwiftCode_CountryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	database.DB = db

	swiftCode := &models.SwiftCode{
		SwiftCode:     "BPHKPLPK",
		IsHeadquarter: false,
		BankName:      "Bank BPH",
		CountryISO2:   "PL",
		CountryName:   "Poland",
		Address:       "ul. Towarowa 25A",
	}

	mock.ExpectExec("^INSERT INTO countries").
		WithArgs(swiftCode.CountryISO2, swiftCode.CountryName).
		WillReturnError(errors.New("database error"))

	_, err = InsertSwiftCode(swiftCode)
	if err == nil {
		t.Error("Expected error but got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestInsertSwiftCode_SwiftCodeError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	database.DB = db

	swiftCode := &models.SwiftCode{
		SwiftCode:     "BPHKPLPK",
		IsHeadquarter: false,
		BankName:      "Bank BPH",
		CountryISO2:   "PL",
		CountryName:   "Poland",
		Address:       "ul. Towarowa 25A",
	}

	mock.ExpectExec("^INSERT INTO countries").
		WithArgs(swiftCode.CountryISO2, swiftCode.CountryName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("^INSERT INTO banks").
		WithArgs(swiftCode.SwiftCode, swiftCode.IsHeadquarter, swiftCode.BankName, swiftCode.Address, swiftCode.CountryISO2).
		WillReturnError(errors.New("database error"))

	_, err = InsertSwiftCode(swiftCode)
	if err == nil {
		t.Error("Expected error but got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetBankCodeAndBranchesBySwift(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	database.DB = db

	swiftCode := "CHASUS33XXX"
	expectedResponse := &models.SwiftCodeResponse{
		SwiftCode:     "CHASUS33XXX",
		BankName:      "JP Morgan Chase",
		Address:       "270 Park Avenue",
		IsHeadquarter: true,
		CountryISO2:   "US",
		CountryName:   "United States",
		Branches: []models.CountriesSwiftCodeResponse{
			{
				SwiftCode:     "CHASUS33NYC",
				BankName:      "JP Morgan Chase NYC",
				Address:       "123 Broadway",
				IsHeadquarter: false,
				CountryISO2:   "US",
			},
		},
	}

	mainRows := sqlmock.NewRows([]string{"swift_code", "bank_name", "address", "is_headquarter", "country_iso2", "country_name"}).
		AddRow(swiftCode, "JP Morgan Chase", "270 Park Avenue", true, "US", "United States")
	mock.ExpectQuery("SELECT").
		WithArgs(swiftCode).
		WillReturnRows(mainRows)

	branchRows := sqlmock.NewRows([]string{"swift_code", "bank_name", "address", "is_headquarter", "country_iso2", "country_name"}).
		AddRow("CHASUS33XXX", "JP Morgan Chase", "270 Park Avenue", true, "US", "United States").
		AddRow("CHASUS33NYC", "JP Morgan Chase NYC", "123 Broadway", false, "US", "United States")
	mock.ExpectQuery("SELECT").
		WithArgs("CHASUS33%").
		WillReturnRows(branchRows)

	result, err := GetBankCodeAndBranchesBySwift(swiftCode)
	if err != nil {
		t.Errorf("Error not expected: %v", err)
	}

	if !reflect.DeepEqual(result, expectedResponse) {
		t.Errorf("Results do not match.\nExpected: %+v\nGot: %+v", expectedResponse, result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetBankCodeAndBranchesBySwift_NonHeadquarter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	database.DB = db

	swiftCode := "BPHKPLPK"
	expectedResponse := &models.SwiftCodeResponse{
		SwiftCode:     "BPHKPLPK",
		BankName:      "Bank BPH",
		Address:       "ul. Towarowa 25A",
		IsHeadquarter: false,
		CountryISO2:   "PL",
		CountryName:   "Poland",
		Branches:      nil,
	}

	rows := sqlmock.NewRows([]string{"swift_code", "bank_name", "address", "is_headquarter", "country_iso2", "country_name"}).
		AddRow(swiftCode, "Bank BPH", "ul. Towarowa 25A", false, "PL", "Poland")
	mock.ExpectQuery("SELECT").
		WithArgs(swiftCode).
		WillReturnRows(rows)

	result, err := GetBankCodeAndBranchesBySwift(swiftCode)
	if err != nil {
		t.Errorf("Error not expected: %v", err)
	}

	if !reflect.DeepEqual(result, expectedResponse) {
		t.Errorf("Results do not match.\nExpected: %+v\nGot: %+v", expectedResponse, result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetBanksByISO2_CountryNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	database.DB = db

	iso2 := "XX"

	mock.ExpectQuery("SELECT iso2, name FROM countries").
		WithArgs(iso2).
		WillReturnRows(sqlmock.NewRows([]string{"iso2", "name"}))

	_, err = GetBanksByISO2(iso2)
	if err == nil {
		t.Error("Expected error but got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
func TestDeleteSwiftCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	database.DB = db

	swiftCode := "BPHKPLPK"

	mock.ExpectExec("DELETE FROM banks").
		WithArgs(swiftCode).
		WillReturnResult(sqlmock.NewResult(0, 1))

	rowsAffected, err := DeleteSwiftCode(swiftCode)
	if err != nil {
		t.Errorf("Error not expected: %v", err)
	}

	if rowsAffected != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowsAffected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteSwiftCode_NoRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	database.DB = db

	swiftCode := "INVALIDCODE"

	mock.ExpectExec("DELETE FROM banks").
		WithArgs(swiftCode).
		WillReturnResult(sqlmock.NewResult(0, 0))

	rowsAffected, err := DeleteSwiftCode(swiftCode)
	if err != nil {
		t.Errorf("Did not expect error, but got: %v", err)
	}

	if rowsAffected != 0 {
		t.Errorf("Expected 0 rows affected, got %d", rowsAffected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
