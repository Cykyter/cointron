package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gopkg.in/telegram-bot-api.v4"
)

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

func main() {
	telegramAPIKey := os.Getenv("CoinTronTelegramAPIKey")
	log.Printf(telegramAPIKey)
	bot, err := tgbotapi.NewBotAPI(telegramAPIKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		messageStr, msgErr := messageHandler(update.Message.Text, bot.Self.UserName)
		if msgErr == nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageStr)
			bot.Send(msg)
		}
	}
}

func enabledBotCommand(s string, botname string) bool {
	str := strings.Split(s, " ")
	if str[0] == fmt.Sprintf("/polo@%s", botname) || str[0] == "/polo" {
		return true
	}
	return false
}

func messageHandler(s string, botname string) (string, error) {
	str := strings.Split(s, " ")
	if str[0] == fmt.Sprintf("/polo@%s", botname) || str[0] == "/polo" {
		if len(str) == 1 {
			// Default to BTC LTC
			data, _ := getCurrentPoloniex("BTC LTC")
			return data, nil
		}

		if len(str) == 3 {
			data, _ := getCurrentPoloniex(fmt.Sprintf("%s %s", strings.ToUpper(str[1]), strings.ToUpper(str[2])))
			return data, nil
		}
	}
	return "", errors.New("Invalid command")
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
	if tickerErr == false {
		return "", errors.New("Invalid ticker value")
	}

	return fmt.Sprintf("Currency: %s\nTime: %s\nLast value: %s\nLowest Ask: %s\nHighest Bid: %s\nSpread: %s\nPercent Change:  %s\nBase Volume: %s\nQuote Volume: %s\n24h High: %s", args, currentTime, ticker.Last, ticker.LowestAsk, ticker.HighestBid, ticker.LowestAsk.Sub(ticker.HighestBid), ticker.PercentChange, ticker.BaseVolume, ticker.QuoteVolume, ticker.High24Hr), nil
}
