package main

import (
	"math/rand"
	"sync"
)

type StockData struct {
	Symbol string
	Price  float64
}

// / FetchStockData fetches the stock data for the given symbols
func FetchStockData(symbols []string) (chan *StockData, error) {
	// Create a channel to receive the stock data
	stockDataChannel := make(chan *StockData, len(symbols))
	// create a wait group to wait for all the go routines to finish
	var wg sync.WaitGroup

	for _, symbol := range symbols {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			// Fetch the stock data for the given symbol
			stockDataChannel <- &StockData{Symbol: symbol, Price: rand.Float64() * 100}
		}(symbol)
	}
	go func() {
		wg.Wait()
		close(stockDataChannel)
	}()
	return stockDataChannel, nil
}
