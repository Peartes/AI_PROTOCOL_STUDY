package main_test

import (
	"testing"

	app "github.com/peartes/scrimba/stock_broker"
	"github.com/stretchr/testify/require"
)

func TestRunApp(t *testing.T) {
	err := app.RunApp()
	require.NoError(t, err)
}
