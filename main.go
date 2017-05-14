package main

import (
	"fmt"
	"log"

	"time"

	"github.com/shopspring/decimal"
	"gopkg.in/jcelliott/turnpike.v2"
	"gopkg.in/telegram-bot-api.v4"
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
	IsFrozen      bool
	High24Hr      decimal.Decimal
	Low24Hr       decimal.Decimal
	Spread        decimal.Decimal
}

// GetPoloniexTicker transforms []interface{} received from poloniex
// websocket api into PoloniexTicker type
func GetPoloniexTicker(args []interface{}) *PoloniexTicker {
	ticker := new(PoloniexTicker)
	ticker.Currency = args[0].(string)
	ticker.Time = time.Now()
	ticker.Last, _ = decimal.NewFromString(args[1].(string))
	ticker.LowestAsk, _ = decimal.NewFromString(args[2].(string))
	ticker.HighestBid, _ = decimal.NewFromString(args[3].(string))
	ticker.PercentChange, _ = decimal.NewFromString(args[4].(string))
	ticker.BaseVolume, _ = decimal.NewFromString(args[5].(string))
	ticker.QuoteVolume, _ = decimal.NewFromString(args[6].(string))
	if args[7].(float64) != 1 {
		ticker.IsFrozen = true
	} else {
		ticker.IsFrozen = false
	}
	ticker.High24Hr, _ = decimal.NewFromString(args[8].(string))
	ticker.Low24Hr, _ = decimal.NewFromString(args[9].(string))
	ticker.Spread = ticker.LowestAsk.Sub(ticker.HighestBid)
	return ticker
}

func main() {
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	fmt.Println("Testing Poloniex WAMP")
	c, err := turnpike.NewWebsocketClient(turnpike.JSON, "wss://api.poloniex.com", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to poloniex router")
	_, err = c.JoinRealm("realm1", nil)

	btcltcData := new(PoloniexTicker)
	onPoloniexTicker := func(args []interface{}, kwargs map[string]interface{}) {
		ticker := GetPoloniexTicker(args)
		if ticker.Currency == "BTC_LTC" {
			log.Println(ticker)
			btcltcData = ticker
		}
	}

	if err := c.Subscribe("ticker", nil, onPoloniexTicker); err != nil {
		log.Fatalln("Error subscribing to ticker:", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == fmt.Sprintf("/polo@%s", bot.Self.UserName) {
			if btcltcData != nil {
				data := btcltcData
				messageStr := fmt.Sprintf("Currency: %s\nTime: %s\nLast value: %s\nLowest Ask: %s\nHighest Bid: %s\nSpread: %s\nPercent Change:  %s\nBase Volume: %s\nQuote Volume: %s\n24h High: %s", data.Currency, data.Time, data.Last, data.LowestAsk, data.HighestBid, data.Spread, data.PercentChange, data.BaseVolume, data.QuoteVolume, data.High24Hr)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageStr)
				bot.Send(msg)
			}
		}
	}

	log.Println("listening for meta events")
	<-c.ReceiveDone
	log.Println("disconnected")
}
