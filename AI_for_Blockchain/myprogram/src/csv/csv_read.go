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

type CSVRecord struct {
	SepalLength float64
	SepalWidth  float64
	PetalLength float64
	PetalWidth  float64
	Species     string
	ParseError  error
}

func ReadCsvFile(filePath string, config *CSVReaderConfig) ([]CSVRecord, error) {
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
	var data []CSVRecord

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading line", err)
			continue
		}
		var csvRecord CSVRecord
		for idx, value := range line {
			// last column is the species
			if idx == 4 {
				if value == "" {
					// species is empty
					csvRecord.ParseError = fmt.Errorf("species is empty")
					break
				} else {
					csvRecord.Species = value
					continue
				}
			}
			var floatValue float64

			_, err := fmt.Sscanf(value, "%f", &floatValue)
			if err != nil {
				csvRecord.ParseError = fmt.Errorf("error parsing float value: %v", err)
				break
			}
			switch idx {
			case 0:
				csvRecord.SepalLength = floatValue
			case 1:
				csvRecord.SepalWidth = floatValue
			case 2:
				csvRecord.PetalLength = floatValue
			case 3:
				csvRecord.PetalWidth = floatValue
			}
		}
		if csvRecord.ParseError == nil {
			data = append(data, csvRecord)
		}
	}

	return data, nil
}
