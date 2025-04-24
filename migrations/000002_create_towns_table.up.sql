CREATE TABLE towns (
    name VARCHAR(50) PRIMARY KEY,
    timezone VARCHAR(40) NOT NULL,
    country VARCHAR(2) NOT NULL,
    CONSTRAINT fk_country FOREIGN KEY (country)
        REFERENCES countries (iso2)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);