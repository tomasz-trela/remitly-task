package models

type CountryResponse struct {
	CountryISO2 string                       `json:"countryISO2"`
	CountryName string                       `json:"countryName"`
	SwiftCodes  []CountriesSwiftCodeResponse `json:"swiftCodes"`
}
