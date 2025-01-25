package main_test

import (
	"testing"

	broker "github.com/peartes/scrimba/stock_broker"
	"github.com/stretchr/testify/require"
)

func TestFetchStockData(t *testing.T) {
	tickers := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA"}

	channel, err := broker.FetchStockData(tickers)

	require.NoError(t, err)
	require.NotNil(t, channel)

	var stockData []broker.PolygonAPIStockData
	for data := <-channel; data != nil; data = <-channel {
		stockData = append(stockData, *data)
	}
	require.Equal(t, len(tickers), len(stockData))
}
