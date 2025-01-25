package main

type StockData struct {
	Symbol string
	Price  float64
}

// StockData represents the top-level structure of the JSON response
type PolygonAPIStockData struct {
	Adjusted     bool                  `json:"adjusted"`
	NextURL      string                `json:"next_url"`
	QueryCount   int                   `json:"queryCount"`
	RequestID    string                `json:"request_id"`
	Results      []PolygonAPIStockInfo `json:"results"`
	ResultsCount int                   `json:"resultsCount"`
	Status       string                `json:"status"`
	Ticker       string                `json:"ticker"`
}

// StockInfo represents individual stock result data
type PolygonAPIStockInfo struct {
	Close          float64 `json:"c"`  // "c" - Closing price
	High           float64 `json:"h"`  // "h" - Highest price
	Low            float64 `json:"l"`  // "l" - Lowest price
	Transactions   int     `json:"n"`  // "n" - Number of trades
	Open           float64 `json:"o"`  // "o" - Opening price
	Timestamp      int64   `json:"t"`  // "t" - Timestamp
	Volume         float64 `json:"v"`  // "v" - Volume traded
	VolumeWeighted float64 `json:"vw"` // "vw" - Volume weighted average price
}
