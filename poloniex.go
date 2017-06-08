package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Constants
const (
	PoloniexBaseURL   = "https://poloniex.com/"
	PoloniexPublicURL = PoloniexBaseURL + "public"
	PoloniexTickerAPI = PoloniexPublicURL + "?command=returnTicker"
)

// PoloniexTicker struct
type PoloniexTicker struct {
	Currency      string
	Time          time.Time
	Last          decimal.Decimal
	LowestAsk     decimal.Decimal
	HighestBid    decimal.Decimal
	PercentChange decimal.Decimal
	BaseVolume    decimal.Decimal
	QuoteVolume   decimal.Decimal
	IsFrozen      string
	High24Hr      decimal.Decimal
	Low24Hr       decimal.Decimal
}

func getCurrentPoloniex(args string) (string, error) {
	currencyPair := strings.Split(args, " ")

	if len(currencyPair) != 2 {
		return "", errors.New("Invalid args")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Get(PoloniexTickerAPI)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	currentTime := time.Now()
	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return "", readErr
	}
	var data map[string]PoloniexTicker
	decodeErr := json.Unmarshal(body, &data)
	if decodeErr != nil {
		return "", decodeErr
	}

	ticker, tickerErr := data[fmt.Sprintf("%s_%s", currencyPair[0], currencyPair[1])]
	if !tickerErr {
		return "", errors.New("Invalid ticker value")
	}

	return fmt.Sprintf("Poloniex\nCurrency: %s\nTime: %s\nLast value: %s\nLowest Ask: %s\nHighest Bid: %s\nSpread: %s\nPercent Change:  %s\nBase Volume: %s\nQuote Volume: %s\n24h High: %s", args, currentTime, ticker.Last, ticker.LowestAsk, ticker.HighestBid, ticker.LowestAsk.Sub(ticker.HighestBid), ticker.PercentChange, ticker.BaseVolume, ticker.QuoteVolume, ticker.High24Hr), nil
}
