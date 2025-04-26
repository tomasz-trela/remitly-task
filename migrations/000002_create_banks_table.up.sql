CREATE TABLE banks (
    swift VARCHAR(50) PRIMARY KEY,
    code VARCHAR(15),
    name TEXT NOT NULL,
    address TEXT,
    country_iso2 VARCHAR(2) NOT NULL,
    CONSTRAINT fk_countries_iso2_banks FOREIGN KEY (country_iso2) REFERENCES countries(iso2) ON DELETE RESTRICT ON UPDATE CASCADE   
);