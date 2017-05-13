package main

import (
	"github.com/shopspring/decimal"
	"fmt"
	"gopkg.in/jcelliott/turnpike.v2"
	"log"
	"gopkg.in/telegram-bot-api.v4"
)

type PoloniexTicker struct {
	Currency      string
	Last          decimal.Decimal
	LowestAsk     decimal.Decimal
	HighestBid    decimal.Decimal
	PercentChange decimal.Decimal
	BaseVolume    decimal.Decimal
	QuoteVolume   decimal.Decimal
	IsFrozen      bool
	High24Hr      decimal.Decimal
	Low24Hr       decimal.Decimal
}

func GetPoloniexTicker(args []interface{}) *PoloniexTicker{
	ticker := new(PoloniexTicker)
	ticker.Currency = args[0].(string)
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
	return ticker
}

func main() {
	bot, err := tgbotapi.NewBotAPI("apikey")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	fmt.Println("Testing Poloniex WAMP")
	c, err := turnpike.NewWebsocketClient(turnpike.JSON, "wss://api.poloniex.com", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to poloniex router")
	_, err = c.JoinRealm("realm1", nil)

	btcData := make(chan *PoloniexTicker)
	if err := c.Subscribe("ticker", nil, func(args []interface{}, kwargs map[string]interface{}) {
		ticker := GetPoloniexTicker(args)
		if ticker.Currency == "BTC_LTC" {
			//log.Println(ticker)
			btcData <- ticker
		}
	}); err != nil {
		log.Fatalln("Error subscribing to ticker:", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 6000
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/polo" {
			if btcData != nil {
				data := <-btcData
				messageStr := fmt.Sprintf("Currency: %s\nLast value: %s\nLowest Ask: %s\nHighest Bid: %s\nPercent Change:  %s\nBase Volume: %s\nQuote Volume: %s\n24h High: %s", data.Currency, data.Last, data.LowestAsk, data.HighestBid, data.PercentChange, data.BaseVolume, data.QuoteVolume, data.High24Hr)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageStr)
				bot.Send(msg)
			}
		}
	}

	log.Println("listening for meta events")
	<-c.ReceiveDone
	log.Println("disconnected")
}

