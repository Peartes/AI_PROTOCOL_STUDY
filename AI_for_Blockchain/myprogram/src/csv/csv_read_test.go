package csv_test

import (
	"myprogram/src/csv"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCSVRead(t *testing.T) {
	// Code for the main function
	fp, err := os.Getwd()
	require.NoError(t, err)

	fileData, err := csv.ReadCsvFile(path.Join(fp, "../", "../", "testdata", "csv", "iris", "iris.data"), &csv.CSVReaderConfig{FieldsPerRecord: 5})
	require.NoError(t, err)

	require.GreaterOrEqual(t, len(fileData), 1)

	for i, row := range fileData {
		require.IsType(t, "string", row.Species, "Row %d is not of the right type string: ", i)
		require.IsType(t, 1.0, row.PetalLength, "Row %d is not of the right type float: ", i)
		require.IsType(t, 1.0, row.PetalWidth, "Row %d is not of the right type float: ", i)
		require.IsType(t, 1.0, row.SepalLength, "Row %d is not of the right type float: ", i)
		require.IsType(t, 1.0, row.SepalWidth, "Row %d is not of the right type float: ", i)
		require.Nil(t, row.ParseError, "Row %d has a parse error: %v", i, row.ParseError)
	}
}
