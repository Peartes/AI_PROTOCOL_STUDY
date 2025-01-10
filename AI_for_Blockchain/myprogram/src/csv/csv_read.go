package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
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

func ReadCsvWithDataFrame(filePath string, config *CSVReaderConfig) (*dataframe.DataFrame, error) {
	irisFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer irisFile.Close()

	irisDF := dataframe.ReadCSV(irisFile)

	return &irisDF, nil
}

// / This filter takes a dataframe and a threshold value and returns
// / only the rows where the sum of the of the SepalLength column and the SepalWidth column is greater than the threshold.
func DataFrameFilterByThreshold(threshold float32, df *dataframe.DataFrame) dataframe.DataFrame {
	// transform the dataframe to a new one
	newDf := df.Rapply(func(ser series.Series) series.Series {
		newSeries := ser.Slice(0, ser.Len()-1)
		sum := func(arr []float64) float64 {
			var total float64
			for _, v := range arr {
				total += v
			}
			return total
		}(newSeries.Float())
		if sum >= float64(threshold) {
			return ser
		}
		ser.Set([]int{0}, series.New([]string{"0"}, series.Float, "temp-series"))
		return ser
	})
	// Filter the dataframe
	return newDf.Filter(dataframe.F{
		Colidx:     0,
		Comparator: "==",
		Comparando: "0.0", // not the cleanest way to do this but it works ans Sepal length is always positive
	})
}
