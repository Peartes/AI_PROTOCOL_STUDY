package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/peartes/scrimba/stock_broker/config"
)

func getStockDataFromPolygonAPI(ticker string, startDate, endDate string, apiKey string) (*PolygonAPIStockData, error) {
	url := fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/%s/range/1/day/%s/%s?apiKey=%s", ticker, startDate, endDate, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data for %s: %v", ticker, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response for %s: %v", ticker, err)
	}

	if resp.StatusCode == 200 {
		var res PolygonAPIStockData
		err := json.Unmarshal(body, &res)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling response for %s: %v", ticker, err)
		}
		return &res, nil
	} else {
		return nil, fmt.Errorf("error fetching stock data for %s: %s", ticker, resp.Status)
	}
}

// / FetchStockData fetches the stock data for the given symbols
func FetchStockData(symbols []string) (chan *PolygonAPIStockData, error) {
	// Create a channel to receive the stock data
	stockDataChannel := make(chan *PolygonAPIStockData, len(symbols))
	// create a wait group to wait for all the go routines to finish
	var wg sync.WaitGroup

	for _, symbol := range symbols {
		wg.Add(1)
		go func(symbol string) {
			apiKey := config.GetPolygonAPIKey()
			defer wg.Done()
			// Fetch the stock data for the given symbol
			res, err := getStockDataFromPolygonAPI(symbol, "2023-01-09", "2023-02-10", apiKey)
			if err != nil {
				fmt.Printf("error getting stock data for symbol %s %v", symbol, err)
			}
			stockDataChannel <- res
		}(symbol)
	}
	go func() {
		wg.Wait()
		close(stockDataChannel)
	}()
	return stockDataChannel, nil
}
