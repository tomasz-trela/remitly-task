package queries

const UpsertCountry = `
INSERT INTO countries (iso2, name)
VALUES ($1, $2)
ON CONFLICT (iso2)
DO UPDATE SET name = EXCLUDED.name
`

const UpsertBank = `
INSERT INTO banks (swift, is_headquarter, name, address, country_iso2)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (swift)
DO UPDATE SET 
  is_headquarter = EXCLUDED.is_headquarter,
  name = EXCLUDED.name,
  address = EXCLUDED.address,
  country_iso2 = EXCLUDED.country_iso2
`

const SelectBanksBySwift = `
SELECT b.swift, b.name AS bank_name, b.address, b.is_headquarter, c.iso2, c.name AS country_name
FROM banks b
JOIN countries c ON b.country_iso2 = c.iso2
WHERE b.swift LIKE $1
`

const SelectCountriesByISO2 = `
SELECT iso2, name
FROM countries
WHERE iso2 LIKE $1
LIMIT 1
`

const SelectBanksByISO2 = `
SELECT b.swift, b.name AS bank_name, b.address, b.is_headquarter
FROM banks b
JOIN countries c ON b.country_iso2 = c.iso2
WHERE iso2 = $1
`

const InsertCountryOrDoNothing = `
INSERT INTO countries (iso2, name)
VALUES ($1, $2)
ON CONFLICT (iso2) DO NOTHING;
`

const InsertSwiftCodeOrDoNothing = `
INSERT INTO banks (swift, is_headquarter, name, address, country_iso2)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (swift) DO NOTHING
`

const DeleteSwiftCode = `
DELETE FROM banks
WHERE swift = $1;
`
