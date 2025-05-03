# SWIFT Code API Service

This is a Go-based application that parses, stores, and serves SWIFT code (BIC) data via a REST API. It is designed to transition SWIFT code data from csv format to a fast, queryable web service.

## 📘 Overview

SWIFT (BIC) codes uniquely identify bank branches worldwide. This application:

- Parses SWIFT data from a given file.
- Identifies and links headquarter vs. branch codes.
- Stores the processed data in a Postgres database.
- Provides RESTful endpoints to manage and retrieve SWIFT information.

## ⚙️ Setup

### 1. Clone the Repository
```bash
git clone https://github.com/tomasz-trela/remitly-task.git
cd remitly-task
```

### 2. Add `.env` file

Create a `.env` file in the project root with the following content:

```env
DATABASE_URL=postgres://postgres:postgres@db:5432/remitly?sslmode=disable
POSTGRES_DB=remitly
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
```

### 3. Build docker compose app
```bash
docker compose build
```

## 🚀 Running the Application

### 1. Run containers
```bash
docker compose up
```

## 🔗 API Endpoints

### 1. Get SWIFT Code Details

**GET** `/v1/swift-codes/{swift-code}`

**Headquarter Response Example:**

```json
{
  "address": "string",
  "bankName": "string",
  "countryISO2": "string",
  "countryName": "string",
  "isHeadquarter": true,
  "swiftCode": "string",
  "branches": [
    {
      "address": "string",
      "bankName": "string",
      "countryISO2": "string",
      "isHeadquarter": false,
      "swiftCode": "string"
    }
  ]
}
```

**Branch Response Example:**
```json
{
  "address": "string",
  "bankName": "string",
  "countryISO2": "string",
  "countryName": "string",
  "isHeadquarter": false,
  "swiftCode": "string"
}
```

## 2. Get All SWIFT Codes by Country

**Endpoint:**  
`GET /v1/swift-codes/country/{countryISO2code}`

**Description:**  
Returns all SWIFT codes (both headquarters and branches) for a specified country using its ISO2 code.

**Response Example:**
```json
{
  "countryISO2": "string",
  "countryName": "string",
  "swiftCodes": [
    {
      "address": "string",
      "bankName": "string",
      "countryISO2": "string",
      "isHeadquarter": true,
      "swiftCode": "string"
    },
    {
      "address": "string",
      "bankName": "string",
      "countryISO2": "string",
      "isHeadquarter": false,
      "swiftCode": "string"
    }
  ]
}
```

## 3. Add a New SWIFT Code

**Endpoint:**  
`POST /v1/swift-codes`

**Description:**  
Adds a new SWIFT code entry to the database. Can be either a headquarter or a branch.

**Request Body Example:**
```json
{
  "address": "string",
  "bankName": "string",
  "countryISO2": "string",
  "countryName": "string",
  "isHeadquarter": true,
  "swiftCode": "string"
}
```

**Response Example:**
```json
{
  "message": "string"
}
```

## 4. Delete a SWIFT Code

**Endpoint:**  
`DELETE /v1/swift-codes/{swift-code}`

**Description:**  
Deletes a SWIFT code entry from the database if it matches the provided code.

**Response Example:**
```json
{
  "message": "string"
}
```