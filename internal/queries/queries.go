package queries

const UpsertCountry = `
INSERT INTO countries (iso2, name)
VALUES ($1, $2)
ON CONFLICT (iso2)
DO UPDATE SET name = EXCLUDED.name
`

const UpsertBank = `
INSERT INTO banks (swift, code, name, address, country_iso2)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (swift)
DO UPDATE SET 
  code = EXCLUDED.code,
  name = EXCLUDED.name,
  address = EXCLUDED.address,
  country_iso2 = EXCLUDED.country_iso2
`
