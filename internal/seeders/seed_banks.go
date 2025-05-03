package seeders

import (
	"fmt"

	"github.com/tomasz-trela/remitly-task/config"
	"github.com/tomasz-trela/remitly-task/internal/parsers"
	"github.com/tomasz-trela/remitly-task/internal/repository"
)

func SeedBanks() {
	records, err := parsers.LoadSwiftRecords(config.BanksCsvPath)
	if err != nil {
		fmt.Println("Error loading records:", err)
		return
	}
	for _, record := range *records {
		_, err := repository.InsertSwiftCode(&record)
		if err != nil {
			fmt.Println("Error inserting record:", err)
		}
	}

}
