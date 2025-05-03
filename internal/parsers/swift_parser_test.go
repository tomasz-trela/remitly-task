package parsers

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tomasz-trela/remitly-task/internal/models"
)

func TestLoadSwiftRecords(t *testing.T) {
	tests := []struct {
		name          string
		fileContent   string
		expectedCodes []models.SwiftCode
		expectError   bool
	}{
		{
			name: "Valid CSV content",
			fileContent: `PL,BPHKPLPK,1234,Bank BPH,ul. Towarowa 25A,Warszawa,Poland
US,CHASUS33XXX,5678,JP Morgan Chase,270 Park Avenue,New York,United States
DE,DEUTDEFF,9012,Deutsche Bank,Taunusanlage 12,Frankfurt am Main,Germany`,
			expectedCodes: []models.SwiftCode{
				{
					SwiftCode:     "BPHKPLPK",
					IsHeadquarter: false,
					BankName:      "Bank BPH",
					CountryISO2:   "PL",
					CountryName:   "Poland",
					Address:       "ul. Towarowa 25A",
				},
				{
					SwiftCode:     "CHASUS33XXX",
					IsHeadquarter: true,
					BankName:      "JP Morgan Chase",
					CountryISO2:   "US",
					CountryName:   "United States",
					Address:       "270 Park Avenue",
				},
				{
					SwiftCode:     "DEUTDEFF",
					IsHeadquarter: false,
					BankName:      "Deutsche Bank",
					CountryISO2:   "DE",
					CountryName:   "Germany",
					Address:       "Taunusanlage 12",
				},
			},
			expectError: false,
		},
		{
			name: "Invalid ISO2 code",
			fileContent: `PL,BPHKPLPK,1234,Bank BPH,ul. Towarowa 25A,Warszawa,Poland
USA,CHASUS33,5678,JP Morgan Chase,270 Park Avenue,New York,United States
DE,DEUTDEFF,9012,Deutsche Bank,Taunusanlage 12,Frankfurt am Main,Germany`,
			expectedCodes: []models.SwiftCode{
				{
					SwiftCode:     "BPHKPLPK",
					IsHeadquarter: false,
					BankName:      "Bank BPH",
					CountryISO2:   "PL",
					CountryName:   "Poland",
					Address:       "ul. Towarowa 25A",
				},
				{
					SwiftCode:     "DEUTDEFF",
					IsHeadquarter: false,
					BankName:      "Deutsche Bank",
					CountryISO2:   "DE",
					CountryName:   "Germany",
					Address:       "Taunusanlage 12",
				},
			},
			expectError: false,
		},
		{
			name:        "Lowercase ISO2 code",
			fileContent: `pl,BPHKPLPK,1234,Bank BPH,ul. Towarowa 25A,Warszawa,Poland`,
			expectedCodes: []models.SwiftCode{
				{
					SwiftCode:     "BPHKPLPK",
					IsHeadquarter: false,
					BankName:      "Bank BPH",
					CountryISO2:   "PL",
					CountryName:   "Poland",
					Address:       "ul. Towarowa 25A",
				},
			},
			expectError: false,
		},
		{
			name:          "Empty file",
			fileContent:   "",
			expectedCodes: []models.SwiftCode{},
			expectError:   false,
		},
		{
			name:        "Address with spaces",
			fileContent: `PL,BPHKPLPK,1234,Bank BPH, ul. Towarowa 25A ,Warszawa,Poland`,
			expectedCodes: []models.SwiftCode{
				{
					SwiftCode:     "BPHKPLPK",
					IsHeadquarter: false,
					BankName:      "Bank BPH",
					CountryISO2:   "PL",
					CountryName:   "Poland",
					Address:       "ul. Towarowa 25A",
				},
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "testfile.csv")

			err := os.WriteFile(tmpFile, []byte(tc.fileContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			records, err := LoadSwiftRecords(tmpFile)

			if tc.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if err == nil {
				if records == nil {
					t.Errorf("Expected non-nil records, got nil")
				} else if tc.name == "Empty file" {

					if len(*records) != 0 {
						t.Errorf("Expected empty records for empty file, got %d records", len(*records))
					}
				} else if !reflect.DeepEqual(*records, tc.expectedCodes) {
					t.Errorf("Records don't match expected.\nGot: %+v\nExpected: %+v", *records, tc.expectedCodes)
				}
			}
		})
	}
}

func TestLoadSwiftRecords_FileNotFound(t *testing.T) {
	_, err := LoadSwiftRecords("non_existent_file.csv")
	if err == nil {
		t.Errorf("Expected error for non-existent file, got nil")
	}
}

func TestLoadSwiftRecords_MalformedCSV(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "malformed.csv")

	content := `PL,BPHKPLPK,1234,Bank BPH,Address,City,Poland
US,CHASUS33,5678,JP Morgan Chase,270 Park Avenue,New York,USA,Extra Field`

	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	records, err := LoadSwiftRecords(tmpFile)

	if err != nil {
		t.Errorf("Expected no error for malformed CSV (should skip bad rows), got %v", err)
		return
	}
	if records == nil {
		t.Errorf("Expected non-nil records, got nil")
	} else if len(*records) != 1 {
		t.Errorf("Expected 1 record, got %d", len(*records))
	} else {
		expected := models.SwiftCode{
			SwiftCode:     "BPHKPLPK",
			IsHeadquarter: false,
			BankName:      "Bank BPH",
			CountryISO2:   "PL",
			CountryName:   "Poland",
			Address:       "Address",
		}

		if !reflect.DeepEqual((*records)[0], expected) {
			t.Errorf("Record doesn't match expected.\nGot: %+v\nExpected: %+v", (*records)[0], expected)
		}
	}
}
