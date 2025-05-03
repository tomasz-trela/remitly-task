package repository

import (
	"fmt"

	"github.com/tomasz-trela/remitly-task/internal/database"
	"github.com/tomasz-trela/remitly-task/internal/models"
	"github.com/tomasz-trela/remitly-task/internal/queries"
)

func UpsertSwiftCode(swiftCode *models.SwiftCode) error {
	_, err := database.DB.Exec(
		queries.InsertCountryOrDoNothing,
		swiftCode.CountryISO2,
		swiftCode.CountryName,
	)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(
		queries.InsertSwiftCodeOrDoNothing,
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
			swiftCode[:8]+"%",
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var branch models.CountriesSwiftCodeResponse
			if err := rows.Scan(
				&branch.SwiftCode,
				&branch.BankName,
				&branch.Address,
				&branch.IsHeadquarter,
				&branch.CountryISO2,
				new(any),
			); err != nil {
				return nil, err
			}

			if branch.SwiftCode != swiftCode {
				swiftCodeResponse.Branches = append(swiftCodeResponse.Branches, branch)
			}
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

func InsertSwiftCode(swiftCode *models.SwiftCode) (int64, error) {
	_, err := database.DB.Exec(
		queries.InsertCountryOrDoNothing,
		swiftCode.CountryISO2,
		swiftCode.CountryName,
	)
	if err != nil {
		return 0, err
	}

	res, err := database.DB.Exec(
		queries.InsertSwiftCodeOrDoNothing,
		swiftCode.SwiftCode,
		swiftCode.IsHeadquarter,
		swiftCode.BankName,
		swiftCode.Address,
		swiftCode.CountryISO2,
	)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func DeleteSwiftCode(swiftCode string) (int64, error) {
	res, err := database.DB.Exec(
		queries.DeleteSwiftCode,
		swiftCode,
	)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
