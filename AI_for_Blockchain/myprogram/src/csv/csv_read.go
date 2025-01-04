package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type CSVReaderConfig struct {
	FieldsPerRecord int
}

func ReadCsvFile(filePath string, config *CSVReaderConfig) ([][]string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file
	reader := csv.NewReader(file)
	if config != nil {
		reader.FieldsPerRecord = config.FieldsPerRecord
	}
	var data [][]string

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading line", err)
			continue
		}
		data = append(data, line)
	}

	return data, nil
}
