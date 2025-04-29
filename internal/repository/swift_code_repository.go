package repository

import (
	"fmt"
	"strings"

	"github.com/tomasz-trela/remitly-task/internal/database"
	"github.com/tomasz-trela/remitly-task/internal/models"
	"github.com/tomasz-trela/remitly-task/internal/queries"
)

func UpsertSwiftCode(swiftCode *models.SwiftCode) error {
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
		swiftCode.IsHeadquarter,
		swiftCode.BankName,
		swiftCode.Address,
		swiftCode.CountryISO2,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetBankCodeAndBranchesBySwift(swiftCode string) (*models.SwiftCodeResponse, error) {
	var swiftCodeResponse models.SwiftCodeResponse
	err := database.DB.QueryRow(
		queries.SelectBanksBySwift,
		swiftCode,
	).Scan(
		&swiftCodeResponse.SwiftCode,
		&swiftCodeResponse.BankName,
		&swiftCodeResponse.Address,
		&swiftCodeResponse.IsHeadquarter,
		&swiftCodeResponse.CountryISO2,
		&swiftCodeResponse.CountryName,
	)
	if err != nil {
		return nil, err
	}

	if swiftCodeResponse.IsHeadquarter {
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
				&branch.IsHeadquarter,
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

func GetBanksByISO2(iso2 string) (*models.CountryResponse, error) {
	var country models.CountryResponse
	err := database.DB.QueryRow(
		queries.SelectCountriesByISO2,
		iso2,
	).Scan(
		&country.CountryISO2,
		&country.CountryName,
	)
	if err != nil {
		return nil, err
	}
	rows, err := database.DB.Query(
		queries.SelectBanksByISO2,
		iso2,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var bank models.CountriesSwiftCodeResponse
		if err := rows.Scan(
			&bank.SwiftCode,
			&bank.BankName,
			&bank.Address,
			&bank.IsHeadquarter,
		); err != nil {
			return nil, err
		}
		bank.CountryISO2 = country.CountryISO2
		country.SwiftCodes = append(country.SwiftCodes, bank)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &country, nil
}

func InsertSwiftCode(swiftCode *models.SwiftCode) error {
	_, err := database.DB.Exec(
		queries.InsertCountryOrDoNothing,
		swiftCode.CountryISO2,
		swiftCode.CountryName,
	)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(
		queries.UpsertBank,
		swiftCode.SwiftCode,
		swiftCode.IsHeadquarter,
		swiftCode.BankName,
		swiftCode.Address,
		swiftCode.CountryISO2,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteSwiftCode(swiftCode string) error {
	_, err := database.DB.Exec(
		queries.DeleteSwiftCode,
		swiftCode,
	)
	if err != nil {
		return err
	}

	return nil
}
