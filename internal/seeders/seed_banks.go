package seeders

import (
	"fmt"

	"github.com/tomasz-trela/remitly-task/internal/parsers"
	"github.com/tomasz-trela/remitly-task/internal/repository"
)

func SeedBanks() {
	records, err := parsers.LoadSwiftRecords()
	if err != nil {
		fmt.Println("Error loading records:", err)
		return
	}
	for _, record := range *records {
		err := repository.UpsertSwiftCode(&record)
		if err != nil {
			fmt.Println("Error inserting record:", err)
		}
	}

}
