CREATE TABLE banks (
    swift VARCHAR(50) PRIMARY KEY,
    code VARCHAR(255),
    name TEXT NOT NULL,
    address TEXT,
    town VARCHAR(255) NOT NULL,
    CONSTRAINT fk_towns_name_banks FOREIGN KEY (town) REFERENCES towns(name)
);