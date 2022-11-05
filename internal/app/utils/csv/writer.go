package csv

import (
	"balance-microservice/internal/app/store"
	"encoding/csv"
	"os"
)

func WriteCVS(records [][]string, filename string) *store.Error {

	csvFile, err := os.Create("./files/csv/" + filename + ".csv")
	if err != nil {
		return &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	csvWriter.Comma = ';'

	for _, record := range records {
		_ = csvWriter.Write(record)
	}

	return nil
}
