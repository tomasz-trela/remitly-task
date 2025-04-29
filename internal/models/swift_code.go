package models

type SwiftCodeResponse struct {
	Address       string                       `json:"address"`
	BankName      string                       `json:"bankName"`
	CountryISO2   string                       `json:"countryISO2"`
	CountryName   string                       `json:"countryName"`
	IsHeadquarter bool                         `json:"isHeadquarter"`
	SwiftCode     string                       `json:"swiftCode"`
	Branches      []CountriesSwiftCodeResponse `json:"branches,omitempty"`
}

type SwiftCode struct {
	SwiftCode     string `json:"swiftCode"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	Address       string `json:"address"`
}

type CountriesSwiftCodeResponse struct {
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
	Address       string `json:"address"`
}
