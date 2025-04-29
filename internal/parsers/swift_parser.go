package parsers

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tomasz-trela/remitly-task/config"
	"github.com/tomasz-trela/remitly-task/internal/models"
)

func LoadSwiftRecords() (*[]models.SwiftCode, error) {
	file, err := os.Open(config.BanksCsvPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var records []models.SwiftCode

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading line:", err)
			return nil, err
		}

		records = append(records, models.SwiftCode{
			SwiftCode:     record[1],
			IsHeadquarter: strings.HasSuffix(record[1], "XXX"),
			BankName:      record[3],
			CountryISO2:   record[0],
			CountryName:   record[6],
			Address:       strings.Trim(record[4], " "),
		})
	}

	return &records, nil
}
