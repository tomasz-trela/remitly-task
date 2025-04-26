package repository

import (
	"fmt"
	"strings"

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

func GetBankCodeBySwift(swiftCode string) (*models.SwiftCodeResponse, error) {
	var swiftCodeResponse models.SwiftCodeResponse
	err := database.DB.QueryRow(
		queries.SelectBanksBySwift,
		swiftCode,
	).Scan(
		&swiftCodeResponse.SwiftCode,
		&swiftCodeResponse.BankName,
		&swiftCodeResponse.Address,
		&swiftCodeResponse.CountryISO2,
		&swiftCodeResponse.CountryName,
	)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(swiftCode, "XXX") {
		swiftCodeResponse.IsHeadquarter = true
		rows, err := database.DB.Query(
			queries.SelectBanksBySwift,
			strings.TrimSuffix(swiftCode, "XXX")+"%",
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var branch models.SwiftCodeResponse
			if err := rows.Scan(
				&branch.SwiftCode,
				&branch.BankName,
				&branch.Address,
				&branch.CountryISO2,
				&branch.CountryName,
			); err != nil {
				return nil, err
			}
			swiftCodeResponse.Branches = append(swiftCodeResponse.Branches, branch)
		}

	} else {
		swiftCodeResponse.IsHeadquarter = false
	}

	fmt.Println("swiftCodeResponse", swiftCodeResponse)

	return &swiftCodeResponse, nil
}
