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
		require.Equal(t, 5, len(row), "Row %d does not have 5 fields: %v", i, row)
	}
}
