package repository

import (
	"github.com/tomasz-trela/remitly-task/internal/database"
	"github.com/tomasz-trela/remitly-task/internal/models"
	"github.com/tomasz-trela/remitly-task/internal/queries"
)

func CreateSwiftCode(swiftCode *models.SwiftCode) error {

	_, err := database.DB.Exec(
		queries.UpsertCountry,
		swiftCode.CountryISO2,
		swiftCode.CountryName,
	)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(
		queries.UpsertBank,
		swiftCode.SwiftCode,
		swiftCode.CodeType,
		swiftCode.BankName,
		swiftCode.Address,
		swiftCode.CountryISO2,
	)
	if err != nil {
		return err
	}

	return nil
}
