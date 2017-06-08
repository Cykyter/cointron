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
	BitfinexBaseURL   = "https://api.bitfinex.com/"
	BitfinexPublicURL = BitfinexBaseURL + "v1/"
	BitfinexTickerAPI = BitfinexPublicURL + "pubticker/"
)

// BitfinexTicker struct
type BitfinexTicker struct {
	CurrencyPair string
	Mid          decimal.Decimal `json:"mid"`
	Bid          decimal.Decimal `json:"bid"`
	Ask          decimal.Decimal `json:"ask"`
	LastPrice    decimal.Decimal `json:"last_price"`
	Low          decimal.Decimal `json:"low"`
	High         decimal.Decimal `json:"high"`
	Volume       decimal.Decimal `json:"volume"`
	Timestamp    string          `json:"timestamp"`
}

func getCurrentBitfinex(args string) (string, error) {
	currencyPair := strings.Split(args, " ")

	if len(currencyPair) != 2 {
		return "", errors.New("Invalid args")
	}

	currencyCommand := currencyPair[0] + currencyPair[1]
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Get(BitfinexTickerAPI + currencyCommand)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	currentTime := time.Now()
	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return "", readErr
	}
	data := new(BitfinexTicker)
	decodeErr := json.Unmarshal(body, &data)
	if decodeErr != nil {
		return "", decodeErr
	}

	if len(data.Timestamp) != 0 {
		return fmt.Sprintf("Bitfinex\nCurrency: %s\nTime: %s\nLast value: %s\nLowest Ask: %s\nHighest Bid: %s\nSpread: %s\nVolume: %s\n24h High: %s\n24h Low: %s", args, currentTime, data.LastPrice, data.Ask, data.Bid, data.Ask.Sub(data.Bid), data.Volume, data.High, data.Low), nil
	}
	return "", errors.New("Invalid currency pair")
}
