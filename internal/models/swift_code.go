package models

type SwiftCodeResponse struct {
	BankName      string              `json:"bankName"`
	CountryISO2   string              `json:"countryISO2"`
	CountryName   string              `json:"countryName"`
	IsHeadquarter bool                `json:"isHeadquarter"`
	SwiftCode     string              `json:"swiftCode"`
	Address       string              `json:"address"`
	Branches      []SwiftCodeResponse `json:"branches,omitempty"`
}

type SwiftCode struct {
	SwiftCode     string
	IsHeadquarter bool
	BankName      string
	CountryISO2   string
	CountryName   string
	Address       string
}

type CountriesSwiftCodeResponse struct {
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
	Address       string `json:"address"`
}
